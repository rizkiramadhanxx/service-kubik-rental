package category

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	var input CreateCategoryRequest
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

	category := entity.Category{Name: input.Name}
	if err := config.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
		"data":    category,
	})
}

func GetAllCategories(c *fiber.Ctx) error {
	var categories []entity.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

func GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var category entity.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Category not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category found",
		"data":    category,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Field %s %s", e.Field(), e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	category.Name = input.Name
	if err := config.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
		"data":    category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Category{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Category deleted successfully"})
}
