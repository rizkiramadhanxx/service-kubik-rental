package device

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Get("/", GetAllDevices)
	router.Post("/multiple-ping", PingMultipleDevices)
	router.Get("/ping/:id", PingDevice)
	router.Get("/action/:id", ActionToDevice)
	router.Post("/", CreateDevice)
	router.Patch("/:id", UpdateDevice)
	router.Delete("/:id", DeleteDevice)
	router.Get("/:id", GetDevice) // <- pindahkan ini ke paling bawah
}
