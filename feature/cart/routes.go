package cart

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("pos"),
	)

	route.Get("/", GetAllCarts)
	route.Get("/:id", GetCartByID)
	route.Post("/", CreateCart)
	route.Patch("/:id", UpdateCart)
	route.Delete("/:id", DeleteCart)
	route.Post("/add-item", AddCartItem)
	route.Post("/toggle-paid/:id", ToggleCartItemPaid)
	route.Post("/action-qty-cart-item", UpdateCartItemQty)
	route.Delete("/cart-item/:id", DeleteCartItem)

}
