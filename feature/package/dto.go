package package_bill

type CreatePackageRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Duration int    `json:"duration" validate:"required,gte=0"`
	Price    int    `json:"price" validate:"required,gt=0"`
}

type UpdatePackageRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Duration int    `json:"duration" validate:"required,gte=0"`
	Price    int    `json:"price" validate:"required,gt=0"`
}

type PackageResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Price     int    `json:"price"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
