package entity

type Package struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name" validate:"required,min=3"`
	Duration  int     `json:"duration" validate:"required,gte=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
