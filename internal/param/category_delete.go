package param

type DeleteCategoryResponse struct {
	Success bool `json:"success"`
}

type DeleteCategoryRequest struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
}