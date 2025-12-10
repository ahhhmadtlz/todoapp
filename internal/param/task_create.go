package param

import "time"

type CreateTaskRequest struct {
	UserID uint `json:"user_id"`
	CategoryID  uint `json:"category_id"`
	Title       string `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority string  `json:"priority,omitempty"` //default  medium
	Status string `json:"status,omitempty"` //default pending
}



type CreateTaskResponse struct {
	Task TaskInfo `json:"task"`
}

