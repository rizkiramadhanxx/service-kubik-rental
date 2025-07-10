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
	"kubik-rental/feature/cart"
	"kubik-rental/feature/category"
	"kubik-rental/feature/device"
	"kubik-rental/feature/member"
	package_bill "kubik-rental/feature/package"
	"kubik-rental/feature/product"
	"kubik-rental/feature/role"
	"kubik-rental/feature/user"
	"kubik-rental/scheduler"
)

func RunMigration() {
	err := config.DB.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Product{},
		&entity.Category{},
		&entity.CartItem{},
		&entity.Device{},
		&entity.Package{},
		&entity.Cart{},
		&entity.Billing{},
		&entity.Member{},
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
	if err := config.SeedAdmin(); err != nil {
		log.Fatal("Seeding admin failed:", err)
	}

	// routes
	user.SetupRoutes(app.Group("/user"))
	role.SetupRoutes(app.Group("/role"))
	product.SetupRoutes(app.Group("/product"))
	category.SetupRoutes(app.Group("/category"))
	auth.SetupRoutes(app.Group("/auth"))
	device.SetupRoutes(app.Group("/device"))
	member.SetupRoutes(app.Group("/member"))
	package_bill.SetupRoutes(app.Group("/package"))
	cart.SetupRoutes(app.Group("/cart"))

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
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
