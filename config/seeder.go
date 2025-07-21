package config

import (
	"encoding/json"
	"fmt"
	"kubik-rental/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin() error {
	fmt.Println("Seeding Admin + Packages + Devices + Menus...")

	// Seed Role & Admin User
	var userCount int64
	DB.Model(&entity.User{}).Count(&userCount)
	if userCount == 0 {
		var role entity.Role
		err := DB.Where("name = ?", "admin").First(&role).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
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

		// Cek jika user "admin" belum ada
		var existing entity.User
		if err := DB.Where("username = ?", "admin").First(&existing).Error; err != nil {
			hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			admin := entity.User{
				Name:     "Admin",
				Username: "admin",
				Password: string(hashed),
				RoleID:   role.ID,
			}
			if err := DB.Create(&admin).Error; err != nil {
				return err
			}
		}
	}

	now := time.Now()

	// Seed Packages (selang-seling PS2/PS3)
	var packageCount int64
	DB.Model(&entity.Package{}).Count(&packageCount)
	if packageCount == 0 {
		packages := []entity.Package{
			{Name: "30 Menit PS2", Duration: 30, Price: 3000, IsLoss: false},
			{Name: "60 Menit PS3", Duration: 60, Price: 6000, IsLoss: false},
			{Name: "90 Menit PS2", Duration: 90, Price: 9000, IsLoss: false},
			{Name: "120 Menit PS3", Duration: 120, Price: 12000, IsLoss: false},
			{Name: "180 Menit PS2", Duration: 180, Price: 18000, IsLoss: true},
		}
		for i := range packages {
			packages[i].CreatedAt = now
			packages[i].UpdatedAt = now
		}
		if err := DB.Create(&packages).Error; err != nil {
			return err
		}
	}

	// Seed Devices (TV 1 - 5)
	var deviceCount int64
	DB.Model(&entity.Device{}).Count(&deviceCount)
	if deviceCount == 0 {
		var devices []entity.Device
		for i := 1; i <= 5; i++ {
			devices = append(devices, entity.Device{
				IP:        fmt.Sprintf("192.168.0.%d", 100+i),
				Name:      fmt.Sprintf("TV %d", i),
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
		if err := DB.Create(&devices).Error; err != nil {
			return err
		}
	}

	// Seed Products (Menus)
	var productCount int64
	DB.Model(&entity.Product{}).Count(&productCount)
	if productCount == 0 {
		products := []entity.Product{
			{Name: "Indomie", Price: 5000},
			{Name: "Gooday", Price: 3000},
			{Name: "Kapal Api", Price: 4000},
			{Name: "Pop Mie", Price: 6000},
			{Name: "Mizone", Price: 5000},
		}
		for i := range products {
			products[i].CreatedAt = now
			products[i].UpdatedAt = now
		}
		if err := DB.Create(&products).Error; err != nil {
			return err
		}
	}

	return nil
}
