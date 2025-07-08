package entity

import "time"

type Product struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID *uint     `gorm:"default:null" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // ← pointer + bisa null
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
