package category

type CreateCategoryRequest struct {
	// minimal 4
	Name string `json:"name" validate:"required,min=4"`
}

type GetCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateCategoryRequest = CreateCategoryRequest
