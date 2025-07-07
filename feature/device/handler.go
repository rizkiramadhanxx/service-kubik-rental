package device

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func isValidAction(action string) bool {
	for _, a := range dto.AllActions {
		if a == action {
			return true
		}
	}
	return false
}

func GetAllDevices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Hitung total data
	var total int64
	if err := config.DB.Model(&entity.Device{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Ambil data paginated
	var devices []entity.Device
	if err := config.DB.Limit(limit).Offset(offset).Find(&devices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Hitung total halaman (jika limit > 0)
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))

		fmt.Println(totalPages)
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]entity.Device]{
		Status:  fiber.StatusOK,
		Message: "Success get all devices",
		Data:    devices,
		Meta: &dto.Meta{
			Page:      page,
			Limit:     limit,
			Total:     int(total),
			TotalPage: totalPages,
		},
	})
}

func GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device entity.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Device not found"})
	}
	return c.Status(fiber.StatusOK).JSON(dto.Response[entity.Device]{
		Status:  fiber.StatusOK,
		Data:    device,
		Message: "Success get device",
	})
}

func ActionToDevice(c *fiber.Ctx) error {

	id := c.Params("id")
	action := c.Query("action")

	if !isValidAction(action) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid action",
		})
	}

	var device entity.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Device not found",
		})
	}

	err := pkg.ExecuteAdbAction(device.IP, action)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(dto.Response[any]{
			Status:  fiber.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: fmt.Sprintf("Action '%s' sent to device", action),
		Data: map[string]interface{}{
			"id":     device.ID,
			"ip":     device.IP,
			"action": action,
		},
	})
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

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Data:    map[string]interface{}{"isReachable": isReachable, "isAdbConnected": isAdbConnected},
		Message: "Success ping device",
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
