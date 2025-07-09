package entity

import "gorm.io/datatypes"

type CartItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CartID      uint           `json:"cart_id"`
	Type        string         `json:"type"`         // "product", "billing"
	ReferenceID uint           `json:"reference_id"` // id dari product atau package
	Name        string         `json:"name"`
	Quantity    int            `json:"quantity"`
	Price       float64        `json:"price"`
	Subtotal    float64        `json:"subtotal"`
	Metadata    datatypes.JSON `json:"metadata"` // pakai import "gorm.io/datatypes"
}
