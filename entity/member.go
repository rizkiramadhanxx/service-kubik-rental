package entity

import "time"

type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
}
