package cart

import (
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCart(c *fiber.Ctx) error {
	var input entity.Cart
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status":  fiber.StatusBadRequest,
			"Message": "Invalid request body",
		})
	}

	input.Status = "active"
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status":  fiber.StatusBadRequest,
			"Message": "Failed to create cart",
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

	query := config.DB.Model(&entity.Cart{}).Preload("Items")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var carts []entity.Cart
	if err := query.Limit(limit).Offset(offset).Find(&carts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPage,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]entity.Cart]{
		Status:  fiber.StatusOK,
		Data:    carts,
		Message: "Carts retrieved successfully",
		Meta:    &meta,
	})
}

func GetCartByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var cart entity.Cart
	if err := config.DB.Preload("Items").First(&cart, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"Status":  fiber.StatusNotFound,
				"Message": "Cart not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status":  fiber.StatusBadRequest,
			"Message": "Failed to retrieve cart",
			"Errors":  err.Error(),
		})
	}
	return c.JSON(dto.Response[entity.Cart]{
		Status:  fiber.StatusOK,
		Message: "Cart retrieved successfully",
		Data:    cart,
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

	cart.Status = input.Status
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

	if err := config.DB.Delete(&entity.Cart{}, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to delete cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Cart deleted successfully",
	})
}
