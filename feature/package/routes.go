package package_bill

import (
	"kubik-rental/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	route := router.Use(
		middleware.Authentication(),
		middleware.RequireModuleAccess("package"),
	)

	route.Get("/", GetAllPackages)
	route.Get("/:id", GetPackageByID)
	route.Post("/", CreatePackage)
	route.Patch("/:id", UpdatePackage)
	route.Delete("/:id", DeletePackage)
}
