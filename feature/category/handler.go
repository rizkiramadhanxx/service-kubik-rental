package category

import (
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	var input CreateCategoryRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  pkg.FormatValidationError(err),
			"status":  fiber.StatusBadRequest,
		})
	}

	// ✅ Validasi manual
	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	category := entity.Category{Name: input.Name}
	if err := config.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
		"status":  fiber.StatusCreated,
	})
}

func GetAllCategories(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	keyword := c.Query("keyword", "")

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := config.DB.Model(&entity.Category{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var categories []entity.Category
	if err := config.DB.Limit(limit).Offset(offset).Where("name LIKE ?", "%"+keyword+"%").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Mapping ke DTO
	var categoryResponses []GetCategoryResponse
	for _, cat := range categories {
		categoryResponses = append(categoryResponses, GetCategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	// Hitung total pages
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPage,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]GetCategoryResponse]{
		Status:  fiber.StatusOK,
		Data:    categoryResponses,
		Message: "Categories found",
		Meta:    &meta,
	})
}

func GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var category entity.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Category not found"})
	}

	// Mapping ke DTO
	categoryResponse := GetCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[GetCategoryResponse]{
		Status:  fiber.StatusOK,
		Data:    categoryResponse,
		Message: "Category found",
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category entity.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Category not found"})
	}

	var input UpdateCategoryRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// ✅ Validasi manual
	if err := pkg.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	category.Name = input.Name
	if err := config.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
		"status":  fiber.StatusOK,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Category{}, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Category deleted successfully"})
}
