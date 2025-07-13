package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint                `gorm:"primaryKey" json:"id"`
	TransactionCode string              `gorm:"uniqueIndex" json:"transaction_code"`
	CartName        string              `json:"cart_name"`
	MemberID        *uint               `json:"member_id,omitempty"`
	BuyerName       *string             `json:"buyer_name,omitempty"`
	IsMember        bool                `json:"is_member"`
	Total           int                 `json:"total"`
	Type            string              `json:"type"` // "product", "billing", "mixed"
	CreatedAt       time.Time           `json:"created_at"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
}

type TransactionDetail struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TransactionID uint   `json:"transaction_id"`
	ItemType      string `json:"item_type"` // "product" / "billing"

	// 🔀 Kolom gabungan billing dan product
	DeviceName  *string    `json:"device_name,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Duration    *int       `json:"duration,omitempty"`
	PackageName *string    `json:"package_name,omitempty"`

	ProductName  *string `json:"product_name,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
	SKU          *string `json:"sku,omitempty"`

	Price    int `json:"price"`
	Qty      int `json:"qty"`
	Subtotal int `json:"subtotal"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // ⬅️ auto isi saat insert
}

// ✅ Auto-generate transaction_code sebelum insert
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.TransactionCode = fmt.Sprintf("TRX-%d", time.Now().Unix())
	return
}
