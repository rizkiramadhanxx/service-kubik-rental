package billing

import (
	"kubik-rental/config"
	"kubik-rental/entity"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BillingInput struct {
	DeviceID  uint      `json:"device_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

func CreateBilling(c *fiber.Ctx) error {

	var input BillingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid input"})
	}

	// Validasi waktu
	if input.EndTime.Before(input.StartTime) {
		return c.Status(400).JSON(fiber.Map{"message": "End time must be after start time"})
	}

	// Simpan billing
	billing := entity.Billing{
		DeviceID:  input.DeviceID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    "active",
	}

	if err := config.DB.Create(&billing).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to create billing"})
	}

	// Update device jadi tidak tersedia
	if err := config.DB.Model(&entity.Device{}).
		Where("id = ?", input.DeviceID).
		Update("available", false).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to update device availability"})
	}

	return c.JSON(fiber.Map{
		"message": "Billing created",
		"data":    billing,
	})
}
