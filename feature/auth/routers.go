package auth

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Post("/login", Login)
	router.Post("/register", Register)
}
