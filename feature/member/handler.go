package member

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateMember(c *fiber.Ctx) error {
	var input CreateMemberRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Field %s %s", e.Field(), e.Tag())
		}
		return c.Status(400).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	member := entity.Member{
		Name:   input.Name,
		Phone:  input.Phone,
		Points: 0,
	}

	if err := config.DB.Create(&member).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	response := MemberResponse{
		ID:        member.ID,
		Name:      member.Name,
		Phone:     member.Phone,
		Points:    member.Points,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.Status(201).JSON(response)
}

func GetMemberByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var member entity.Member
	if err := config.DB.First(&member, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Member not found"})
	}

	response := MemberResponse{
		ID:        member.ID,
		Name:      member.Name,
		Phone:     member.Phone,
		Points:    member.Points,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(response)
}

func GetAllMembers(c *fiber.Ctx) error {
	var members []entity.Member
	if err := config.DB.Find(&members).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	var responses []MemberResponse
	for _, m := range members {
		responses = append(responses, MemberResponse{
			ID:        m.ID,
			Name:      m.Name,
			Phone:     m.Phone,
			Points:    m.Points,
			CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(responses)
}

func UpdateMember(c *fiber.Ctx) error {
	id := c.Params("id")
	var member entity.Member
	if err := config.DB.First(&member, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Member not found"})
	}

	var input UpdateMemberRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Field %s %s", e.Field(), e.Tag())
		}
		return c.Status(400).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	member.Name = input.Name
	member.Phone = input.Phone

	if err := config.DB.Save(&member).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Member updated successfully"})
}

func DeleteMember(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Member{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Member deleted successfully"})
}
