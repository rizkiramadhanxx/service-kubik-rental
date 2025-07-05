package member

type CreateMemberRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type UpdateMemberRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type MemberResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Points    int    `json:"points"`
	CreatedAt string `json:"created_at"`
}
