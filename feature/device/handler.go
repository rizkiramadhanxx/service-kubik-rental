package device

import (
	"fmt"
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"strconv"
	"time"

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
	keyword := c.Query("keyword", "")
	filterBilling := c.Query("is_billing", "") // bisa "true" atau "false"

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query awal
	query := config.DB.Model(&entity.Device{})

	// Filter by keyword
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR ip LIKE ?", like, like)
	}

	// Hitung total data (sebelum pagination)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Ambil data + preload semua billings (tanpa kondisi)
	var devices []entity.Device
	if err := query.
		Preload("Billings"). // preload semua billing
		Limit(limit).
		Offset(offset).
		Find(&devices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Transform + hitung is_billing
	type DeviceWithFlag struct {
		entity.Device
		IsBilling bool `json:"is_billing"`
	}

	var result []DeviceWithFlag
	for _, device := range devices {
		isBilling := false
		for _, b := range device.Billings {
			if b.EndTime.IsZero() || b.EndTime.After(time.Now()) {
				isBilling = true
				break
			}
		}
		result = append(result, DeviceWithFlag{
			Device:    device,
			IsBilling: isBilling,
		})
	}

	// Filter `is_billing` jika diminta
	if filterBilling != "" {
		wantBilling := filterBilling == "true"
		filtered := []DeviceWithFlag{}
		for _, d := range result {
			if d.IsBilling == wantBilling {
				filtered = append(filtered, d)
			}
		}
		result = filtered
		total = int64(len(result))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]DeviceWithFlag]{
		Status:  fiber.StatusOK,
		Message: "Success get all devices",
		Data:    result,
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
	isAdbConnected := pkg.ConnectAdbToIP(device.IP)

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Data:    map[string]interface{}{"isReachable": isReachable, "isAdbConnected": isAdbConnected},
		Message: "Success ping device",
	})
}

func PingMultipleDevices(c *fiber.Ctx) error {
	var req PingMultipleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := pkg.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   pkg.FormatValidationError(err),
		})
	}

	var results []map[string]interface{}

	for _, id := range req.IDs {
		var device entity.Device
		if err := config.DB.First(&device, id).Error; err != nil {
			results = append(results, map[string]interface{}{
				"id":             id,
				"isReachable":    false,
				"isAdbConnected": false,
				"error":          "Device not found",
			})
			continue
		}

		isReachable := pkg.PingIP(device.IP)
		isAdbConnected := pkg.CheckAdbConnected(device.IP)

		results = append(results, map[string]interface{}{
			"id":             device.ID,
			"ip":             device.IP,
			"name":           device.Name,
			"isReachable":    isReachable,
			"isAdbConnected": isAdbConnected,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{
		Status:  fiber.StatusOK,
		Message: "Success ping multiple devices",
		Data:    results,
	})
}

func CreateDevice(c *fiber.Ctx) error {
	var device CreateDeviceRequest
	if err := c.BodyParser(&device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  pkg.FormatValidationError(err),
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := pkg.Validate.Struct(device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  pkg.FormatValidationError(err),
		})
	}

	newDevice := entity.Device{
		Name: device.Name,
		IP:   device.IP,
	}

	if err := config.DB.Create(&newDevice).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Device created successfully", "status": fiber.StatusCreated})
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

	return c.JSON(fiber.Map{"message": "Device updated successfully", "status": fiber.StatusOK})
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&entity.Device{}, id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response[any]{
			Status:  fiber.StatusInternalServerError,
			Message: result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{
			Status:  fiber.StatusNotFound,
			Message: "Device not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Device deleted successfully", "status": fiber.StatusOK})
}
