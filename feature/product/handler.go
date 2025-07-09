package product

import (
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	var input CreateProductRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  pkg.FormatValidationError(err),
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	// validate category tapi boleh null
	if input.CategoryID != nil {
		var category entity.Category
		if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Category not found"})
		}
	}

	product := entity.Product{
		Name:       input.Name,
		SKU:        input.SKU,
		Price:      input.Price,
		Stock:      input.Stock,
		CategoryID: input.CategoryID,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"status":  fiber.StatusCreated,
	})
}

func GetAllProducts(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	keyword := c.Query("keyword", "")
	category := c.Query("category", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Build base query
	query := config.DB.Model(&entity.Product{}).Preload("Category")

	// Conditionally add filters (mirip "keyword || undefined")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category_id = ?", category)
	}

	// Hitung total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Ambil data
	var products []entity.Product
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Format response
	var productResponses []GetProductResponse
	for _, p := range products {
		productResponses = append(productResponses, GetProductResponse{
			ID:       p.ID,
			Name:     p.Name,
			SKU:      p.SKU,
			Price:    p.Price,
			Stock:    p.Stock,
			Category: p.Category,
		})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPage,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]GetProductResponse]{
		Status:  fiber.StatusOK,
		Data:    productResponses,
		Message: "Products found",
		Meta:    &meta,
	})
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product entity.Product
	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	productResponse := GetProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		SKU:      product.SKU,
		Price:    product.Price,
		Stock:    product.Stock,
		Category: product.Category,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[GetProductResponse]{
		Status:  fiber.StatusOK,
		Data:    productResponse,
		Message: "Product found",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if input.CategoryID != nil {
		var category entity.Category
		if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Category not found"})
		}
	}

	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	product.Name = input.Name
	product.SKU = input.SKU
	product.Price = input.Price
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID

	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"status":  fiber.StatusOK,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&entity.Product{}, id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}
