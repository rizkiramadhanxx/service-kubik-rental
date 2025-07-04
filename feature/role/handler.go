// feature/role/handler.go
package role

import (
	"encoding/json"
	"fmt"
	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func IsValidModule(name string) bool {
	for _, m := range entity.AllModules {
		if m == name {
			return true
		}
	}
	return false
}

func CreateRole(c *fiber.Ctx) error {
	var input entity.Role

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

	// Validate modules value
	var parsedModules []string
	if err := json.Unmarshal([]byte(input.Modules), &parsedModules); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid modules format. Must be JSON array of strings"})
	}
	for _, m := range parsedModules {
		if !IsValidModule(m) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid module: %s", m)})
		}
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(input)
}

func GetAllRoles(c *fiber.Ctx) error {
	var roles []entity.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(roles)
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var role entity.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}
	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role entity.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}

	var input entity.Role
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Validate JSON format and enum
	var parsedModules []string
	if err := json.Unmarshal([]byte(input.Modules), &parsedModules); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid modules format"})
	}
	for _, m := range parsedModules {
		if !IsValidModule(m) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid module: %s", m)})
		}
	}

	role.Name = input.Name
	role.Modules = input.Modules

	if err := config.DB.Save(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(role)
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Role{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Role deleted successfully"})
}
