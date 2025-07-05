package entity

import "time"

type Billing struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeviceID  uint      `json:"device_id"`
	Device    Device    `gorm:"foreignKey:DeviceID" json:"device"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `gorm:"type:varchar(20)" json:"status"` // active, expired
	CreatedAt time.Time `json:"created_at"`
}
