package dto

type CreateMemberDTO struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type UpdateMemberDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
