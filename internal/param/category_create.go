package param

type CreateCategoryRequest struct {
	UserID      uint    `json:"-"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreatecategoryResponse struct {
	Category CategoryInfo `json:"category"`
}
