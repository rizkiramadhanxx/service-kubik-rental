package entity

import "time"

type CartItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CartID     uint      `json:"cart_id" validate:"required"`
	ItemType   string    `json:"item_type" validate:"required,oneof=product billing"`
	ProductID  *uint     `json:"product_id,omitempty"`
	Product    *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"` // ⬅️ Tambahkan ini
	BillingID  *uint     `json:"billing_id,omitempty"`
	Billing    *Billing  `gorm:"foreignKey:BillingID;references:ID" json:"billing,omitempty"` // 🟢
	Duration   *int      `json:"duration,omitempty"`
	Price      int       `json:"price" validate:"required,gt=0"`
	TotalPrice int       `json:"total_price"`
	Qty        int       `json:"qty" validate:"required,min=1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
