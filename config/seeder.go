package config

import (
	"encoding/json"
	"kubik-rental/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin() error {
	var count int64
	DB.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil // skip seeding
	}

	// cek atau buat role admin
	var role entity.Role
	err := DB.Where("name = ?", "admin").First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// JSON stringify: ["user","setting","role","transaction"]
			modulesJSON, _ := json.Marshal(entity.AllModules)

			role = entity.Role{
				Name:    "admin",
				Modules: string(modulesJSON),
			}
			if err := DB.Create(&role).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// buat user admin
	hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := entity.User{
		Name:     "Admin",
		Username: "admin",
		Password: string(hashed),
		RoleID:   role.ID,
	}

	return DB.Create(&admin).Error
}
