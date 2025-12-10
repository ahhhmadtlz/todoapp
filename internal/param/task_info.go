package param

import "time"

type TaskInfo struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	CategoryID  uint       `json:"category_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}


