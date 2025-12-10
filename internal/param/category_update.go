package param

type UpdateCategoryRequest struct {
	ID uint `json:"-"`
	UserID uint `json:"-"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}