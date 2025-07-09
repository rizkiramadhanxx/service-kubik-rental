package entity

import "time"

type Device struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IP        string    `json:"ip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Billing   *Billing  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"billing,omitempty"`
}
