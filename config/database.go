package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("kubik.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Enable PRAGMA foreign_keys = ON
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from GORM DB:", err)
	}

	// Execute PRAGMA to enable foreign key constraints
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}

	DB = db
}
