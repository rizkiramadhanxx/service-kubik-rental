package role

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("role"),
	)

	route.Get("/module", getAllModule)
	route.Get("/", GetAllRoles)
	route.Get("/:id", GetRoleByID)
	route.Post("/", CreateRole)
	route.Patch("/:id", UpdateRole)
	route.Delete("/:id", DeleteRole)
}
