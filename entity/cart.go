package entity

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"cart_items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
