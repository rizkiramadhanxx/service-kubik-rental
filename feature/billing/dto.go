package billing

type CreateBillingWithOptionalCartRequest struct {
	CartID    *uint   `json:"cart_id,omitempty"`
	CartName  *string `json:"cart_name,omitempty"` // Jika tidak ada cart_id
	DeviceID  uint    `json:"device_id" validate:"required"`
	PackageID uint    `json:"package_id" validate:"required"`
}
