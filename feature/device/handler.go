package device

import (
	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/pkg"

	"github.com/gofiber/fiber/v2"
)

func GetAllDevices(c *fiber.Ctx) error {
	var devices []entity.Device
	if err := config.DB.Find(&devices).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(devices)
}

func GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device entity.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Device not found"})
	}
	return c.JSON(device)
}

func PingDevice(c *fiber.Ctx) error {
	id := c.Params("id")

	// Ambil device dari DB
	var device entity.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Device not found"})
	}

	// Panggil fungsi dari pkg
	isReachable := pkg.PingIP(device.IP)
	isAdbConnected := pkg.CheckAdbConnected(device.IP)

	return c.JSON(fiber.Map{
		"device":        device,
		"ip_reachable":  isReachable,
		"adb_connected": isAdbConnected,
	})
}

func CreateDevice(c *fiber.Ctx) error {
	var device entity.Device
	if err := c.BodyParser(&device); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	if err := config.DB.Create(&device).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(device)
}

func UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device entity.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Device not found"})
	}

	if err := c.BodyParser(&device); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if err := config.DB.Save(&device).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(device)
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&entity.Device{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(204)
}
