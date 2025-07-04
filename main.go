package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/feature/auth"
	"kubik-rental/feature/role"
	"kubik-rental/feature/user"
)

func RunMigration() {
	err := config.DB.AutoMigrate(
		&entity.User{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
func main() {
	app := fiber.New()
	config.LoadEnv()

	// Allow all origins, methods, headers
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	fmt.Println("Server running on port 3000")

	config.InitDB()
	RunMigration()

	api := app.Group("/api")
	user.SetupRoutes(api.Group("/user"))
	role.SetupRoutes(api.Group("/role"))
	auth.SetupRoutes(api.Group("/auth"))

	app.Listen(":3000")
}
