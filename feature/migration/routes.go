package migration

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router) {
	handler := NewMigrationHandler()

	// Migration routes
	app.Post("/run", handler.RunMigration)         // POST /migration/run
	app.Get("/status", handler.GetMigrationStatus) // GET /migration/status
}
