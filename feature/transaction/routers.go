package transaction

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("transaction"),
	)

	router.Post("/checkout", CheckoutFromCart)
	router.Get("/", GetAllTransaction)
}
