package config

import (
	"log"
	"time" // ⬅️ tambahkan

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// ⬇️ Set time.Local ke Asia/Jakarta agar semua time.Now() dan autoCreateTime pakai WIB
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal("Gagal load lokasi Asia/Jakarta:", err)
	}
	time.Local = loc

	db, err := gorm.Open(sqlite.Open("kubik.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from GORM DB:", err)
	}

	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}

	DB = db
}
