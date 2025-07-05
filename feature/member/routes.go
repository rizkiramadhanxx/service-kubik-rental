package member

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("user"),
	)

	route.Get("/", GetAllMembers)
	route.Get("/:id", GetMemberByID)
	route.Post("/", CreateMember)
	route.Put("/:id", UpdateMember)
	route.Delete("/:id", DeleteMember)
}
