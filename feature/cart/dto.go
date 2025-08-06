package cart

import "kubik-rental/entity"

type CreateCartItemRequest struct {
	CartID    uint   `json:"cart_id" validate:"required"`
	ItemType  string `json:"item_type" validate:"required,oneof=product billing"`
	ProductID *uint  `json:"product_id,omitempty"`
	BillingID *uint  `json:"billing_id,omitempty"`
	Qty       int    `json:"qty" validate:"required,min=1"`
	Duration  *int   `json:"duration,omitempty"` // untuk billing
}

type CartDetailResponse struct {
	entity.Cart
	TotalPrice   int `json:"total_price"`
	TotalProduct int `json:"total_product"`
	TotalBilling int `json:"total_billing"`
	TotalPay     int `json:"total_pay"`
}

type UpdateQtyRequest struct {
	CartID    uint   `json:"cart_id" validate:"required"`
	ItemType  string `json:"item_type" validate:"required,oneof=product billing"`
	ProductID *uint  `json:"product_id,omitempty"`
	BillingID *uint  `json:"billing_id,omitempty"`
	Action    string `json:"action" validate:"required,oneof=increment decrement set"`
	Value     *int   `json:"value,omitempty"` // hanya wajib jika action == set
}
