package transaction

import (
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CheckoutFromCart(c *fiber.Ctx) error {
	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid request"})
	}
	if err := pkg.Validate.Struct(req); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": err.Error(), "errors": pkg.FormatValidationError(err)})
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		var cart entity.Cart
		if err := tx.Preload("CartItems.Product.Category").
			Preload("CartItems.Billing.Package").
			Preload("CartItems.Billing.Device").
			First(&cart, req.CartID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"status": fiber.StatusNotFound, "message": "Cart tidak ditemukan"})
		}

		if len(cart.CartItems) == 0 {
			return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": false, "message": "Cart kosong"})
		}

		var (
			details   []entity.TransactionDetail
			total     int
			typeSet   = map[string]bool{}
			buyerName *string
			isMember  = false
			memberID  *uint
		)

		if req.MemberID != nil {
			var member entity.Member
			if err := tx.First(&member, *req.MemberID).Error; err == nil {
				buyerName = &member.Name
				isMember = true
				memberID = &member.ID
			}
		} else if req.BuyerName != nil && *req.BuyerName != "" {
			buyerName = req.BuyerName
		}

		for _, item := range cart.CartItems {
			subtotal := item.Price * item.Qty
			total += subtotal
			typeSet[item.ItemType] = true

			td := entity.TransactionDetail{
				ItemType: item.ItemType,
				Qty:      item.Qty,
				Price:    item.Price,
				Subtotal: subtotal,
			}

			switch item.ItemType {
			case "product":
				if item.Product != nil {
					td.ProductName = &item.Product.Name
					td.SKU = item.Product.SKU
					if item.Product.Category != nil {
						td.CategoryName = &item.Product.Category.Name
					}
				}
			case "billing":
				if item.Billing != nil {
					td.DeviceName = &item.Billing.Device.Name
					td.StartTime = &item.Billing.StartTime
					td.EndTime = &item.Billing.EndTime
					td.Duration = item.Duration
					if item.Billing.Package.ID != 0 {
						td.PackageName = &item.Billing.Package.Name
					}
					// Hapus billing aktif
					if err := tx.Delete(&entity.Billing{}, item.Billing.ID).Error; err != nil {
						return err
					}
				}
			}

			details = append(details, td)
		}

		// Tentukan jenis transaksi
		var txType string
		switch {
		case typeSet["product"] && typeSet["billing"]:
			txType = "mixed"
		case typeSet["product"]:
			txType = "product"
		case typeSet["billing"]:
			txType = "billing"
		default:
			txType = "unknown"
		}

		transaction := entity.Transaction{
			CartName:  cart.Name,
			BuyerName: buyerName,
			IsMember:  isMember,
			MemberID:  memberID,
			Total:     total,
			Type:      txType,
			Details:   details,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entity.Cart{}, cart.ID).Error; err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Checkout berhasil",
			"data":    transaction,
		})
	})
}

func GetAllTransaction(c *fiber.Ctx) error {
	var q TransactionQuery
	if err := c.QueryParser(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query params",
		})
	}

	// Default pagination
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	offset := (q.Page - 1) * q.Limit

	tx := config.DB.Model(&entity.Transaction{})

	// Filter: date range
	if q.StartDate != "" && q.EndDate != "" {
		start, err1 := time.Parse("2006-01-02", q.StartDate)
		end, err2 := time.Parse("2006-01-02", q.EndDate)
		if err1 == nil && err2 == nil {
			end = end.Add(24 * time.Hour)
			tx = tx.Where("created_at BETWEEN ? AND ?", start, end)
		}
	}

	// Filter: valid type only
	switch q.Type {
	case "product", "billing", "mixed":
		tx = tx.Where("type = ?", q.Type)
	}

	// Filter: buyer_name keyword
	if q.Keyword != "" {
		tx = tx.Where("buyer_name LIKE ?", "%"+q.Keyword+"%")
	}

	// Count total
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Gagal menghitung data transaksi",
		})
	}

	// Get data
	var transactions []entity.Transaction
	if err := tx.Order("created_at DESC").Limit(q.Limit).Offset(offset).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Gagal mengambil data transaksi",
		})
	}

	if transactions == nil {
		transactions = []entity.Transaction{}
	}

	meta := &dto.Meta{
		Page:      q.Page,
		Limit:     q.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(q.Limit))),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]entity.Transaction]{
		Status:  fiber.StatusOK,
		Message: "Berhasil mengambil data transaksi",
		Data:    transactions,
		Meta:    meta,
	})
}
