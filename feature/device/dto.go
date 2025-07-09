package device

type CreateDeviceRequest struct {
	IP   string `json:"ip" validate:"required,ip"`
	Name string `json:"name" validate:"required,min=3"`
}

type UpdateDeviceRequest struct {
	IP   string `json:"ip" validate:"required,ip"`
	Name string `json:"name" validate:"required,min=3"`
}

type DeviceResponse struct {
	ID        uint   `json:"id"`
	IP        string `json:"ip"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
