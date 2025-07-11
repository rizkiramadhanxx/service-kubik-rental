package cart

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

func CreateCart(c *fiber.Ctx) error {
	var input entity.Cart
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to create cart",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response[entity.Cart]{
		Status:  fiber.StatusCreated,
		Message: "Cart created successfully",
	})
}

func GetAllCarts(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	status := c.Query("status", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := config.DB.Model(&entity.Cart{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var carts []entity.Cart
	if err := query.
		Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Billing").
		Preload("CartItems.Billing.Package").
		Preload("CartItems.Billing.Device").
		Limit(limit).
		Offset(offset).
		Find(&carts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Hitung total2
	cartResponses := make([]CartDetailResponse, 0) // 🟢 penting agar [] bukan null
	for _, cart := range carts {
		var totalPrice, totalProduct, totalBilling int
		for _, item := range cart.CartItems {
			totalPrice += item.TotalPrice
			switch item.ItemType {
			case "product":
				totalProduct++
			case "billing":
				totalBilling++
			}
		}
		cartResponses = append(cartResponses, CartDetailResponse{
			Cart:         cart,
			TotalPrice:   totalPrice,
			TotalProduct: totalProduct,
			TotalBilling: totalBilling,
		})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPage,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]CartDetailResponse]{
		Status:  fiber.StatusOK,
		Message: "Carts retrieved successfully",
		Data:    cartResponses,
		Meta:    &meta,
	})
}

func GetCartByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var cart entity.Cart
	err := config.DB.Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Billing").
		Preload("CartItems.Billing.Package").
		Preload("CartItems.Billing.Device").
		First(&cart, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Cart not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to retrieve cart",
		})
	}

	// Perhitungan total
	var totalPrice, totalProduct, totalBilling int
	for _, item := range cart.CartItems {
		totalPrice += item.TotalPrice
		if item.ItemType == "product" {
			totalProduct++
		} else if item.ItemType == "billing" {
			totalBilling++
		}
	}

	return c.JSON(dto.Response[CartDetailResponse]{
		Status:  fiber.StatusOK,
		Message: "Cart retrieved successfully 1",
		Data: CartDetailResponse{
			Cart:         cart,
			TotalPrice:   totalPrice,
			TotalProduct: totalProduct,
			TotalBilling: totalBilling,
		},
	})
}

func UpdateCart(c *fiber.Ctx) error {
	id := c.Params("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	var cart entity.Cart
	if err := config.DB.First(&cart, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Cart not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to retrieve cart",
		})
	}

	cart.UpdatedAt = time.Now()
	if err := config.DB.Save(&cart).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to update cart",
		})
	}

	return c.JSON(dto.Response[entity.Cart]{
		Status:  fiber.StatusOK,
		Message: "Cart updated successfully",
		Data:    cart,
	})
}

func DeleteCart(c *fiber.Ctx) error {
	id := c.Params("id")

	// Ambil semua CartItem yang terkait
	var cartItems []entity.CartItem
	if err := config.DB.Where("cart_id = ?", id).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get cart items",
		})
	}

	// Kumpulkan BillingID dari CartItem
	var billingIDs []uint
	for _, item := range cartItems {
		if item.BillingID != nil {
			billingIDs = append(billingIDs, *item.BillingID)
		}
	}

	// Hapus billing terkait
	if len(billingIDs) > 0 {
		if err := config.DB.Delete(&entity.Billing{}, billingIDs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to delete related billing",
			})
		}
	}

	// Hapus cart (otomatis akan hapus cart item karena cascade)
	if err := config.DB.Delete(&entity.Cart{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Cart deleted successfully",
	})
}

func AddCartItem(c *fiber.Ctx) error {
	var input CreateCartItemRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	// Validasi Cart
	var cart entity.Cart
	if err := config.DB.First(&cart, input.CartID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Cart not found",
		})
	}

	// Cek apakah item sudah ada
	var existingItem entity.CartItem
	lookup := config.DB.Where("cart_id = ? AND item_type = ?", input.CartID, input.ItemType)

	if input.ItemType == "product" && input.ProductID != nil {
		lookup = lookup.Where("product_id = ?", input.ProductID)
	} else if input.ItemType == "billing" && input.BillingID != nil {
		lookup = lookup.Where("billing_id = ?", input.BillingID)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid item reference",
		})
	}

	if err := lookup.First(&existingItem).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Item already exists in the cart",
		})
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to check existing cart item",
		})
	}

	// Hitung harga berdasarkan jenis item
	var price int
	var totalPrice int
	if input.ItemType == "product" && input.ProductID != nil {
		var product entity.Product
		if err := config.DB.First(&product, input.ProductID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Product not found",
			})
		}
		price = int(product.Price)
		totalPrice = int(product.Price * float64(input.Qty))
	} else if input.ItemType == "billing" && input.BillingID != nil {
		var pkg entity.Billing
		if err := config.DB.Preload("Package").First(&pkg, input.BillingID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Billing not found",
			})
		}

		price = pkg.Package.Price // diasumsikan total harga langsung
	}

	// Buat item baru
	item := entity.CartItem{
		CartID:     input.CartID,
		ItemType:   input.ItemType,
		ProductID:  input.ProductID,
		BillingID:  input.BillingID,
		Qty:        input.Qty,
		Duration:   input.Duration,
		Price:      price,
		CreatedAt:  time.Now(),
		TotalPrice: totalPrice,
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to create cart item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response[entity.CartItem]{
		Status:  fiber.StatusCreated,
		Message: "Cart item added successfully",
		Data:    item,
	})
}

func UpdateCartItemQty(c *fiber.Ctx) error {

	var input UpdateQtyRequest
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

	// Cari cart item
	var cartItem entity.CartItem
	query := config.DB.Where("cart_id = ? AND item_type = ?", input.CartID, input.ItemType)

	if input.ItemType == "product" && input.ProductID != nil {
		query = query.Where("product_id = ?", input.ProductID)
	} else if input.ItemType == "billing" && input.BillingID != nil {
		query = query.Where("billing_id = ?", input.BillingID)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid item reference",
		})
	}

	err := query.First(&cartItem).Error
	notFound := err == gorm.ErrRecordNotFound

	// Jika tidak ditemukan dan action = increment/set → buat item baru
	if notFound && (input.Action == "increment" || input.Action == "set") {
		var price int
		var totalPrice int
		var duration *int

		if input.ItemType == "product" && input.ProductID != nil {
			var product entity.Product
			if err := config.DB.First(&product, input.ProductID).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
					Status:  fiber.StatusNotFound,
					Message: "Product not found",
				})
			}
			price = int(product.Price)
			qty := 1
			if input.Action == "set" && input.Value != nil {
				qty = *input.Value
			}
			totalPrice = price * qty
			cartItem = entity.CartItem{
				CartID:     input.CartID,
				ItemType:   input.ItemType,
				ProductID:  input.ProductID,
				Qty:        qty,
				Price:      price,
				TotalPrice: totalPrice,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
		} else if input.ItemType == "billing" && input.BillingID != nil {
			var pkg entity.Billing
			if err := config.DB.Preload("Package").First(&pkg, input.BillingID).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
					Status:  fiber.StatusNotFound,
					Message: "Package not found",
				})
			}
			price = pkg.Package.Price
			duration = &pkg.Package.Duration
			cartItem = entity.CartItem{
				CartID:     input.CartID,
				ItemType:   input.ItemType,
				BillingID:  input.BillingID,
				Qty:        1,
				Duration:   duration,
				Price:      price,
				TotalPrice: price,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
		}

		if err := config.DB.Create(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
				Status:  fiber.StatusBadRequest,
				Message: "Failed to create cart item",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(dto.Response[entity.CartItem]{
			Status:  fiber.StatusCreated,
			Message: "Cart item created",
			Data:    cartItem,
		})
	} else if notFound {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Cart item not found",
		})
	}

	// Jika ditemukan → update atau hapus jika qty 0
	switch input.Action {
	case "increment":
		cartItem.Qty++
	case "decrement":
		cartItem.Qty--
		if cartItem.Qty < 1 {
			if err := config.DB.Delete(&cartItem).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
					Status:  fiber.StatusBadRequest,
					Message: "Failed to delete cart item",
				})
			}
			return c.JSON(dto.Response[any]{
				Status:  fiber.StatusOK,
				Message: "Cart item deleted (qty became 0)",
			})
		}
	case "set":
		if input.Value == nil || *input.Value < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
				Status:  fiber.StatusBadRequest,
				Message: "Invalid set value (must be >= 1)",
			})
		}
		cartItem.Qty = *input.Value
	}

	// Update ulang totalPrice
	cartItem.TotalPrice = cartItem.Price * cartItem.Qty
	cartItem.UpdatedAt = time.Now()

	if err := config.DB.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to update cart item",
		})
	}

	return c.JSON(dto.Response[entity.CartItem]{
		Status:  fiber.StatusOK,
		Message: "Cart item updated successfully",
		Data:    cartItem,
	})
}

func DeleteCartItem(c *fiber.Ctx) error {
	id := c.Params("id")

	// Ambil dulu cart item-nya
	var item entity.CartItem
	if err := config.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
				Status:  fiber.StatusNotFound,
				Message: "Cart item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve cart item",
		})
	}

	// Jika item adalah billing dan BillingID tidak null, hapus billing
	if item.ItemType == "billing" && item.BillingID != nil {
		if err := config.DB.Delete(&entity.Billing{}, *item.BillingID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to delete related billing",
			})
		}
	}

	// Hapus cart item-nya
	if err := config.DB.Delete(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete cart item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Cart item deleted successfully",
	})
}
