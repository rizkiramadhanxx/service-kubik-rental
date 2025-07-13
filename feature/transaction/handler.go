package transaction

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CheckoutFromCart(c *fiber.Ctx) error {
	currentUser, ok := c.Locals("user").(entity.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
		})
	}

	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request",
		})
	}

	if err := pkg.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"errors":  pkg.FormatValidationError(err),
		})
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		var cart entity.Cart
		if err := tx.Preload("CartItems.Product.Category").
			Preload("CartItems.Billing.Package").
			Preload("CartItems.Billing.Device").
			First(&cart, req.CartID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "Cart tidak ditemukan",
			})
		}

		if len(cart.CartItems) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Cart kosong",
			})
		}

		var (
			details   []entity.TransactionDetail
			total     int
			typeSet   = map[string]bool{}
			buyerName *string
			isMember  = false
			memberID  *uint
		)

		// Cek apakah member
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

		// Bangun detail transaksi dan proses stok
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
					// Kurangi stok
					res := tx.Model(&entity.Product{}).
						Where("id = ? AND stock >= ?", item.Product.ID, item.Qty).
						UpdateColumn("stock", gorm.Expr("stock - ?", item.Qty))

					if res.Error != nil {
						return res.Error
					}
					if res.RowsAffected == 0 {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status":  false,
							"message": "Stok tidak mencukupi untuk produk: " + item.Product.Name,
						})
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
					// ⛔ Jangan hapus Billing di sini (FK masih aktif)
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
			Cashier:   currentUser.Username,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// ✅ Hapus semua CartItems dulu agar FK tidak dilanggar saat hapus billing
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&entity.CartItem{}).Error; err != nil {
			return err
		}

		// ✅ Hapus Billing setelah CartItem sudah tidak ada
		for _, item := range cart.CartItems {
			if item.Billing != nil {
				if err := tx.Delete(&entity.Billing{}, item.Billing.ID).Error; err != nil {
					return err
				}
			}
		}

		// ✅ Hapus cart
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
		tx = tx.Where("buyer_name IS NOT NULL AND buyer_name LIKE ?", "%"+q.Keyword+"%")
	}

	// Count total
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Gagal menghitung data transaksi",
		})
	}

	// Get data
	var transactions []entity.Transaction
	if err := tx.Order("created_at DESC").Limit(q.Limit).Offset(offset).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
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

func GetDetailTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

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

	// Ambil transaksi
	var tx entity.Transaction
	if err := config.DB.First(&tx, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Transaksi tidak ditemukan",
		})
	}

	// Hitung total detail yang cocok
	detailQuery := config.DB.Model(&entity.TransactionDetail{}).Where("transaction_id = ?", id)

	// Filter by type (item_type: billing/product)
	if q.Type != "" {
		detailQuery = detailQuery.Where("item_type = ?", q.Type)
	}

	// Filter by keyword (device_name atau product_name)
	if q.Keyword != "" {
		keyword := "%" + q.Keyword + "%"
		detailQuery = detailQuery.Where(
			"(device_name LIKE ? OR product_name LIKE ?)", keyword, keyword,
		)
	}

	var total int64
	if err := detailQuery.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Gagal menghitung detail transaksi",
		})
	}

	// Ambil data detail sesuai paginasi
	var details []entity.TransactionDetail
	if err := detailQuery.
		Order("id ASC").
		Limit(q.Limit).
		Offset(offset).
		Find(&details).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Gagal mengambil detail transaksi",
		})
	}

	meta := &dto.Meta{
		Page:      q.Page,
		Limit:     q.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(q.Limit))),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]entity.TransactionDetail]{
		Status:  fiber.StatusOK,
		Message: "Berhasil mengambil detail transaksi",
		Data:    details,
		Meta:    meta,
	})

}

type TransactionStatRow struct {
	Date    string `json:"date"`
	All     int    `json:"all"`
	Product int    `json:"product"`
	Billing int    `json:"billing"`
}

func GetTransactionStats(c *fiber.Ctx) error {
	db := config.DB

	monthly := c.Query("monthly") // format: 07-2025
	yearly := c.Query("yearly")   // format: 2025

	if (monthly != "" && yearly != "") || (monthly == "" && yearly == "") {
		return c.JSON(dto.Response[any]{
			Message: "Harus pilih salah satu: 'monthly' (MM-YYYY) atau 'yearly' (YYYY)",
			Status:  fiber.StatusOK,
			Data: fiber.Map{
				"analytic":      []TransactionStatRow{},
				"total":         0,
				"total_product": 0,
				"total_billing": 0,
			},
		})
	}

	type RawResult struct {
		Date     string
		ItemType string
		Total    int
	}
	type TransactionStatRow struct {
		Date    string `json:"date"`
		All     int    `json:"all"`
		Product int    `json:"product"`
		Billing int    `json:"billing"`
	}

	var (
		start         time.Time
		end           time.Time
		layoutGroupBy string
		results       []TransactionStatRow
		statsMap      = map[string]*TransactionStatRow{}
		rawResults    []RawResult

		totalAll     int
		totalProduct int
		totalBilling int
	)

	if monthly != "" {
		t, err := time.Parse("01-2006", monthly)
		if err != nil {
			return c.JSON(dto.Response[any]{
				Message: "Format 'monthly' salah. Gunakan MM-YYYY",
				Status:  fiber.StatusOK,
				Data: fiber.Map{
					"analytic":      []TransactionStatRow{},
					"total":         0,
					"total_product": 0,
					"total_billing": 0,
				},
			})
		}
		start = t
		end = t.AddDate(0, 1, 0)
		layoutGroupBy = "%Y-%m-%d"

		for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
			key := d.Format("2006-01-02")
			statsMap[key] = &TransactionStatRow{Date: key}
		}
	}

	if yearly != "" {
		t, err := time.Parse("2006", yearly)
		if err != nil {
			return c.JSON(dto.Response[any]{
				Message: "Format 'yearly' salah. Gunakan YYYY",
				Status:  fiber.StatusOK,
				Data: fiber.Map{
					"analytic":      []TransactionStatRow{},
					"total":         0,
					"total_product": 0,
					"total_billing": 0,
				},
			})
		}
		start = t
		end = t.AddDate(1, 0, 0)
		layoutGroupBy = "%Y-%m"

		for m := 1; m <= 12; m++ {
			key := fmt.Sprintf("%d-%02d", t.Year(), m)
			statsMap[key] = &TransactionStatRow{Date: key}
		}
	}

	err := db.Model(&entity.TransactionDetail{}).
		Select(fmt.Sprintf("strftime('%s', created_at) as date, item_type, SUM(subtotal) as total", layoutGroupBy)).
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("date, item_type").
		Order("date").
		Scan(&rawResults).Error

	if err != nil {
		return c.JSON(dto.Response[any]{
			Message: "Gagal mengambil data analitik",
			Status:  fiber.StatusOK,
			Data: fiber.Map{
				"analytic":      []TransactionStatRow{},
				"total":         0,
				"total_product": 0,
				"total_billing": 0,
			},
		})
	}

	for _, r := range rawResults {
		row := statsMap[r.Date]
		if row == nil {
			row = &TransactionStatRow{Date: r.Date}
			statsMap[r.Date] = row
		}
		row.All += r.Total
		totalAll += r.Total

		switch r.ItemType {
		case "product":
			row.Product += r.Total
			totalProduct += r.Total
		case "billing":
			row.Billing += r.Total
			totalBilling += r.Total
		}
	}

	for _, v := range statsMap {
		results = append(results, *v)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date < results[j].Date
	})

	return c.JSON(dto.Response[any]{
		Message: "Berhasil mengambil data statistik transaksi",
		Status:  fiber.StatusOK,
		Data: fiber.Map{
			"analytic":      results,
			"total":         totalAll,
			"total_product": totalProduct,
			"total_billing": totalBilling,
		},
	})
}
