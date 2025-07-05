package product

type CreateProductRequest struct {
	Name       string  `json:"name" validate:"required"`
	SKU        string  `json:"sku" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	Stock      int     `json:"stock" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

type UpdateProductRequest = CreateProductRequest
