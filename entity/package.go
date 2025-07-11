package entity

import "time"

type Package struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Duration  int       `json:"duration" validate:"required,gte=0"`
	Price     int       `json:"price" validate:"required,gt=0"`
	IsLoss    bool      `json:"is_loss"`
	CreatedAt time.Time `json:"created_at"` // ✅ ubah ke time.Time
	UpdatedAt time.Time `json:"updated_at"` // ✅ ubah ke time.Time
}
