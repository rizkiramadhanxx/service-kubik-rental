package user

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("user"),
	)

	router.Get("/", GetAllUsers)
	router.Get("/:id", GetUser)
	router.Post("/", CreateUser)
	router.Patch("/:id", UpdateUser)
	router.Delete("/:id", DeleteUser)
}
