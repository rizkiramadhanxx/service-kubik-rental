package dto

import "kubik-rental/entity"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

type UserResponse struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	RoleID   uint        `json:"role_id"`
	Role     entity.Role `json:"role"`
	Username string      `json:"username"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Username string  `json:"username" validate:"required"`
	RoleID   uint    `json:"role_id" validate:"required"`
	Password *string `json:"password" validate:"omitempty,min=5"` // pointer
}
