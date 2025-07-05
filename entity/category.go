package entity

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	Products []Product `json:"-"` // optional, only if you need reverse relation
}
