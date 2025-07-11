package billing

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("billing"),
	)

	route.Post("/", CreateBillingAndInsertToCartHandler)
	// loss billing
	route.Get("/stop-billing", StopLossBilling)

}
