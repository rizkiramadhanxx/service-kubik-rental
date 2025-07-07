package user

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	var input dto.CreateUserRequest

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

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	// check if role exists
	var role entity.Role
	if err := config.DB.First(&role, input.RoleID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Role not found"})
	}

	newUser := entity.User{
		Name:     input.Name,
		Password: string(hashed),
		RoleID:   input.RoleID,
		Role:     role,
		Username: input.Username,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// preload role
	config.DB.Preload("Role").First(&newUser, newUser.ID)

	response := dto.UserResponse{
		ID:     newUser.ID,
		Name:   newUser.Name,
		RoleID: newUser.RoleID,
		Role:   newUser.Role,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entity.User

	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	response := dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		RoleID:   user.RoleID,
		Role:     user.Role,
		Username: user.Username,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[dto.UserResponse]{Status: fiber.StatusOK, Data: response, Message: "User found"})
}

func GetAllUsers(c *fiber.Ctx) error {
	// Ambil query param
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := config.DB.Model(&entity.User{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to count users",
		})
	}

	var users []entity.User
	if err := config.DB.Preload("Role").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Mapping ke DTO
	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			RoleID:   user.RoleID,
			Role:     user.Role,
			Username: user.Username,
		}
	}

	// Hitung total pages
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	// Meta
	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPage,
	}

	// Final JSON response
	return c.Status(fiber.StatusOK).JSON(dto.Response[[]dto.UserResponse]{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    response,
		Meta:    &meta,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entity.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	var input entity.User
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

	// validate role
	if err := config.DB.First(&entity.Role{}, input.Role).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Role not found"})
	}

	// hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	user.Name = input.Name
	user.Password = string(hashed)
	user.Role = input.Role

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.User{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
