package user

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("user"),
	)

	route.Get("/", GetAllUsers)
	route.Post("/", CreateUser)
	route.Get("/:id", GetUser)
	route.Put("/:id", UpdateUser)
	route.Delete("/:id", DeleteUser)
}
