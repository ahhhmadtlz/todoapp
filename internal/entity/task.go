package entity

import "time"

type Priority string
type Status string



const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)
const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusDone       Status = "done"
)

func (p Priority) String() string {
	return string(p)
}

func (s Status) String() string{
	return string(s)
}

func MapToPriorityEntity(priority string) Priority {
	switch priority {
	case "low":
		return PriorityLow
	case "medium":
		return PriorityMedium
	case "high":
		return PriorityHigh
	default:
		return PriorityLow
	}
}

func MapToStatusEntity(status string) Status {
	switch status {
	case "pending":
		return StatusPending
	case "inprogress":
		return StatusInProgress
	case "done":
		return StatusDone
	default:
		return StatusPending
	}
}

type Task struct {
	ID          uint
	UserID      uint
	CategoryID  uint
	Title       string
	DueDate     *time.Time
	Description string
	Priority    Priority
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
