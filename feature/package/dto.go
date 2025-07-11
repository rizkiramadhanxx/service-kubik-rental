package package_bill

import "time"

type CreatePackageRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Duration int    `json:"duration" validate:"required,gte=0"`
	Price    int    `json:"price" validate:"required,gt=0"`
	IsLoss   bool   `json:"is_loss"`
}

type UpdatePackageRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Duration int    `json:"duration" validate:"required,gte=0"`
	Price    int    `json:"price" validate:"required,gt=0"`
	IsLoss   bool   `json:"is_loss"`
}

type PackageResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Duration  int       `json:"duration"`
	Price     int       `json:"price"`
	IsLoss    bool      `json:"is_loss"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
