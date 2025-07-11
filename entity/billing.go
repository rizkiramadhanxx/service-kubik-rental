package entity

import "time"

type Billing struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeviceID  uint      `json:"device_id"`
	Device    Device    `gorm:"foreignKey:DeviceID" json:"device"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	PackageID uint      `json:"package_id"`
	IsActive  bool      `json:"is_active"`
	Package   Package   `gorm:"foreignKey:PackageID;references:ID"` // 🟢
	CreatedAt time.Time `json:"created_at"`
}
