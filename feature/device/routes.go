package device

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Get("/", GetAllDevices)
	router.Get("/:id", GetDevice)
	router.Post("/", CreateDevice)
	router.Put("/:id", UpdateDevice)
	router.Delete("/:id", DeleteDevice)
}
