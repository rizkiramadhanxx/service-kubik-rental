package product

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	var input CreateProductRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Field %s %s", e.Field(), e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	// validasi kategori
	var category entity.Category
	if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid category"})
	}

	product := entity.Product{
		Name:       input.Name,
		SKU:        input.SKU,
		Price:      input.Price,
		Stock:      input.Stock,
		CategoryID: input.CategoryID,
		CreatedAt:  time.Now(),
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// preload category
	config.DB.Preload("Category").First(&product, product.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"data":    product,
	})
}

func GetAllProducts(c *fiber.Ctx) error {
	var products []entity.Product
	if err := config.DB.Preload("Category").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product list retrieved successfully",
		"data":    products,
	})
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product entity.Product
	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product found",
		"data":    product,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product entity.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	var input UpdateProductRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Field %s %s", e.Field(), e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	var category entity.Category
	if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid category"})
	}

	product.Name = input.Name
	product.SKU = input.SKU
	product.Price = input.Price
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID

	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Product{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
