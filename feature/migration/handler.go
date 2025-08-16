package migration

import (
	"log"
	"net/http"

	"kubik-rental/config"
	"kubik-rental/entity"

	"github.com/gofiber/fiber/v2"
)

type MigrationHandler struct{}

func NewMigrationHandler() *MigrationHandler {
	return &MigrationHandler{}
}

// RunMigration menjalankan migration database
func (h *MigrationHandler) RunMigration(c *fiber.Ctx) error {
	log.Println("Memulai migration database...")

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
		&entity.Transaction{},
		&entity.TransactionDetail{},
	)

	if err != nil {
		log.Printf("Migration gagal: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Migration database gagal",
			"error":   err.Error(),
		})
	}

	log.Println("Migration database berhasil")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Migration database berhasil",
		"data": fiber.Map{
			"tables": []string{
				"users",
				"roles",
				"products",
				"categories",
				"cart_items",
				"devices",
				"packages",
				"carts",
				"billings",
				"members",
				"transactions",
				"transaction_details",
			},
		},
	})
}

// GetMigrationStatus mengecek status migration
func (h *MigrationHandler) GetMigrationStatus(c *fiber.Ctx) error {
	// Cek apakah tabel-tabel sudah ada
	tables := []string{
		"users",
		"roles",
		"products",
		"categories",
		"cart_items",
		"devices",
		"packages",
		"carts",
		"billings",
		"members",
		"transactions",
		"transaction_details",
	}

	existingTables := []string{}
	missingTables := []string{}

	for _, table := range tables {
		if config.DB.Migrator().HasTable(table) {
			existingTables = append(existingTables, table)
		} else {
			missingTables = append(missingTables, table)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_tables":    len(tables),
			"existing_tables": existingTables,
			"missing_tables":  missingTables,
			"is_complete":     len(missingTables) == 0,
		},
	})
}


