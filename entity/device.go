package entity

import "time"

type Device struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    bool      `json:"status"`
	IP        string    `json:"ip"`
	Name      string    `json:"name"`
	Available bool      `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
