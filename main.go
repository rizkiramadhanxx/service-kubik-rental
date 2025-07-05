package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"kubik-rental/config"
	"kubik-rental/entity"
	"kubik-rental/feature/auth"
	"kubik-rental/feature/device"
	"kubik-rental/feature/member"
	"kubik-rental/feature/product"
	"kubik-rental/feature/role"
	"kubik-rental/feature/user"
	"kubik-rental/scheduler"
)

func RunMigration() {
	err := config.DB.AutoMigrate(
		&entity.User{}, &entity.Role{}, &entity.Device{},
		&entity.Billing{},
		&entity.Member{},
		&entity.Product{},
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

	// database
	config.InitDB()
	RunMigration()

	// routes
	api := app.Group("/api")
	user.SetupRoutes(api.Group("/user"))
	role.SetupRoutes(api.Group("/role"))
	auth.SetupRoutes(api.Group("/auth"))
	device.SetupRoutes(api.Group("/device"))
	member.SetupRoutes(api.Group("/member"))
	product.SetupRoutes(api.Group("/product"))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/status-adb", func(c *fiber.Ctx) error {
		adbPath := "embed/platform-tools/adb.exe" // relatif ke working directory

		// jalankan adb devices
		cmd := exec.Command(adbPath, "version")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("ADB error: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"output": string(out),
		})
	})

	// Start billing polling
	scheduler.StartBillingPolling()
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3000"
	}
	fmt.Printf("Listening on port %s\n", PORT)

	app.Listen(":" + PORT)
}
