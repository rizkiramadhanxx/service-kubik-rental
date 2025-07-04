package role

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Get("/", GetAllRoles)
	router.Get("/:id", GetRoleByID)
	router.Post("/", CreateRole)
	router.Put("/:id", UpdateRole)
	router.Delete("/:id", DeleteRole)
}
