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

	// ✅ Validasi device
	var device entity.Device
	if err := config.DB.First(&device, input.DeviceID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Device not found",
		})
	}

	// ✅ Validasi package
	var pkgData entity.Package
	if err := config.DB.First(&pkgData, input.PackageID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Package not found",
		})
	}

	// ✅ Cek apakah device sedang digunakan oleh billing aktif lain
	var count int64
	config.DB.Model(&entity.Billing{}).
		Where("device_id = ? AND is_active = ?", input.DeviceID, true).
		Count(&count)

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Device is already in use by another active billing",
		})
	}

	// Hitung endTime
	var endTime *time.Time
	if !pkgData.IsLoss {
		et := time.Now().Add(time.Duration(pkgData.Duration) * time.Minute)
		endTime = &et
	}

	// Simpan billing
	billing := entity.Billing{
		DeviceID:  input.DeviceID,
		PackageID: input.PackageID,
		IsActive:  true,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
	}
	if endTime != nil {
		billing.EndTime = *endTime
	}
	if err := config.DB.Create(&billing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create billing",
		})
	}

	// Ambil atau buat cart
	var cart entity.Cart
	if input.CartID != nil {
		// ✅ Validasi cart
		if err := config.DB.First(&cart, *input.CartID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Cart not found",
			})
		}
	} else if input.CartName != nil {
		cart = entity.Cart{
			Name:      *input.CartName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := config.DB.Create(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to create new cart",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "cart_id or cart_name must be provided",
		})
	}

	// Tambahkan ke cart_items
	cartItem := entity.CartItem{
		CartID:     cart.ID,
		ItemType:   "billing",
		BillingID:  &billing.ID,
		Qty:        1,
		Duration:   &pkgData.Duration,
		Price:      pkgData.Price,
		TotalPrice: pkgData.Price,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := config.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to add cart item",
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

	// 2. Validasi apakah sudah dihentikan
	if !billing.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Billing already stopped",
		})
	}

	// 3. Hitung durasi dan harga total
	now := time.Now()
	durationMin := int(now.Sub(billing.StartTime).Minutes())
	pricePerHour := billing.Package.Price
	totalPrice := int(math.Ceil(float64(durationMin)/60.0)) * pricePerHour

	// 4. Update billing dengan EndTime
	billing.EndTime = now
	if err := config.DB.Save(&billing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update billing",
		})
	}

	// 5. Update cart_item
	if err := config.DB.Model(&entity.CartItem{}).
		Where("billing_id = ?", billing.ID).
		Updates(map[string]interface{}{
			"duration":    durationMin,
			"price":       pricePerHour,
			"total_price": totalPrice,
			"updated_at":  now,
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update cart item",
		})
	}

	return c.JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Billing stopped and cart item updated successfully",
	})
}
