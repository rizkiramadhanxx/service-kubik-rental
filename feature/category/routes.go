package category

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("category"),
	)

	route.Get("/", GetAllCategories)
	route.Get("/:id", GetCategoryByID)
	route.Post("/", CreateCategory)
	route.Patch("/:id", UpdateCategory)
	route.Delete("/:id", DeleteCategory)
}
