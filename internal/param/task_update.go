package param

import "time"

type UpdateTaskRequest struct {
	ID uint  `json:"task_id"`
	UserID uint `json:"user_id"`
	CategoryID  *uint `json:"category_id,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status *string `json:"status,omitempty"`// "pending" or "complete"
	Priority *string  `json:"priority,omitempty"`
}

type UpdateTaskResponse struct {
	Task TaskInfo `json:"task"`
}

