package product

import "kubik-rental/entity"

type CreateProductRequest struct {
	Name       string  `json:"name" validate:"required,min=3"`
	SKU        *string `json:"sku" validate:"omitempty,min=3"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	Stock      int     `json:"stock" validate:"required,gte=0"`
	CategoryID *uint   `json:"category_id" validate:"omitempty"`
}

type GetProductResponse struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	SKU      *string          `json:"sku"`
	Price    float64          `json:"price"`
	Stock    int              `json:"stock"`
	Category *entity.Category `json:"category,omitempty"`
}

type UpdateProductRequest = CreateProductRequest
