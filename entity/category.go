package entity

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	Products []Product `json:"products"` // optional, only if you need reverse relation
}
