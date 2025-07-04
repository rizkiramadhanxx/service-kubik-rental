package user

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	userGroup := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("user"),
	)

	userGroup.Get("/", GetAllUsers)
	userGroup.Post("/", CreateUser)
	userGroup.Get("/:id", GetUser)
	userGroup.Put("/:id", UpdateUser)
	userGroup.Delete("/:id", DeleteUser)
}
