package entity

import "time"

type Transaction struct {
	ID        uint                `gorm:"primaryKey"`
	CartName  string              `json:"cart_name"`
	MemberID  *uint               `json:"member_id,omitempty"`  // nullable
	BuyerName *string             `json:"buyer_name,omitempty"` // snapshot nama pembeli
	IsMember  bool                `json:"is_member"`            // true jika member
	Total     int                 `json:"total"`
	Type      string              `json:"type"` // "product", "billing", "mixed"
	CreatedAt time.Time           `json:"created_at"`
	Details   []TransactionDetail `gorm:"foreignKey:TransactionID"`
}

type TransactionDetail struct {
	ID            uint   `gorm:"primaryKey"`
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
}
