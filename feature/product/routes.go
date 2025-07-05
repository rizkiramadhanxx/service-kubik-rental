package product

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("product"),
	)

	route.Get("/", GetAllProducts)
	route.Get("/:id", GetProductByID)
	route.Post("/", CreateProduct)
	route.Put("/:id", UpdateProduct)
	route.Delete("/:id", DeleteProduct)
}
