package category

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest = CreateCategoryRequest
