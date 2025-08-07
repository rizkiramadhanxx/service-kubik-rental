package billing

import (
	"errors"
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateBillingAndInsertToCartHandler(c *fiber.Ctx) error {
	var input CreateBillingWithOptionalCartRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}
	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Validation error",
			Data:    pkg.FormatValidationError(err),
		})
	}

	var cartItem entity.CartItem

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// ✅ Validasi device
		var device entity.Device
		if err := tx.First(&device, input.DeviceID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Device not found")
		}

		// ✅ Validasi package
		var pkgData entity.Package
		if err := tx.First(&pkgData, input.PackageID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Package not found")
		}

		// ✅ Cek apakah device sedang digunakan oleh billing aktif lain
		var count int64
		tx.Model(&entity.Billing{}).
			Where("device_id = ? AND is_active = ?", input.DeviceID, true).
			Count(&count)
		if count > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Device is already in use by another active billing")
		}

		// Tentukan endTime & harga
		var endTime *time.Time
		price := 0
		totalPrice := 0

		if input.IsLoss != nil && *input.IsLoss {
			// loss mode → tidak ada endTime & harga 0
			endTime = nil
		} else {
			// normal
			et := time.Now().Add(time.Duration(pkgData.Duration) * time.Minute)
			endTime = &et
			price = pkgData.Price
			totalPrice = pkgData.Price
		}

		// Simpan billing
		billing := entity.Billing{
			DeviceID:  input.DeviceID,
			PackageID: input.PackageID,
			IsActive:  true,
			IsLoss:    *input.IsLoss,
			StartTime: time.Now(),
			CreatedAt: time.Now(),
			EndTime:   endTime,
		}

		if err := tx.Create(&billing).Error; err != nil {
			return err
		}

		// Ambil atau buat cart
		var cart entity.Cart
		if input.CartID != nil {
			if err := tx.First(&cart, *input.CartID).Error; err != nil {
				return fiber.NewError(fiber.StatusNotFound, "Cart not found")
			}
		} else if input.CartName != nil {
			cart = entity.Cart{
				Name:      *input.CartName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := tx.Create(&cart).Error; err != nil {
				return err
			}
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "cart_id or cart_name must be provided")
		}

		// Tambahkan cart item
		cartItem = entity.CartItem{
			CartID:     cart.ID,
			ItemType:   "billing",
			BillingID:  &billing.ID,
			Qty:        1,
			Price:      price,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := tx.Create(&cartItem).Error; err != nil {
			return err
		}

		return nil // ✅ commit
	})

	if err != nil {
		// Tangkap error dari dalam transaction dan kirim response
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		return c.Status(code).JSON(dto.Response[any]{
			Status:  code,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response[entity.CartItem]{
		Status:  fiber.StatusCreated,
		Message: "Billing created and added to cart successfully",
		Data:    cartItem,
	})
}

func StopLossBilling(c *fiber.Ctx) error {
	id := c.Params("id")

	// 1. Ambil billing + preload package
	var billing entity.Billing
	if err := config.DB.Preload("Package").First(&billing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Billing not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve billing",
		})
	}

	// 2. Validasi package tidak nil
	if billing.Package.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Billing does not have a valid package",
		})
	}

	// 3. Validasi apakah sudah dihentikan (pointer-safe)
	if billing.EndTime != nil && !billing.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Billing already stopped",
		})
	}

	// 4. Hitung durasi dan harga total (per menit)
	now := time.Now()
	durationMin := int(now.Sub(billing.StartTime).Minutes())
	if durationMin <= 0 {
		durationMin = 1 // minimal 1 menit
	}

	pricePerHour := billing.Package.Price
	totalPrice := int(math.Ceil(float64(durationMin) * float64(pricePerHour) / 60.0))

	// 5. Transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 5a. Update billing lengkap
		billing.EndTime = &now
		billing.IsActive = false

		if err := tx.Save(&billing).Error; err != nil {
			return err
		}

		// 5b. Update cart_item (jika ada)
		if err := tx.Model(&entity.CartItem{}).
			Where("billing_id = ?", billing.ID).
			Updates(map[string]interface{}{
				"price":       totalPrice,
				"total_price": totalPrice,
				"updated_at":  now,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to stop billing transaction",
		})
	}

	// 6. Response lengkap
	return c.JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Billing stopped and cart item updated successfully",
		Data: map[string]any{
			"duration_min": durationMin,
			"price_hour":   pricePerHour,
			"total_price":  totalPrice,
		},
	})
}
