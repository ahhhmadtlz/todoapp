package taskservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) GetTaskByID(ctx context.Context, id uint, userID uint) (param.GetTaskResponse, error) {
	const op = richerror.Op("taskservice.GetTaskByID")

	task, err := s.repo.GetTaskByID(ctx, id, userID)
	if err != nil {
		return param.GetTaskResponse{}, richerror.New(op).WithErr(err)
	}

	return param.GetTaskResponse{
		Task: param.TaskInfo{
			ID:          task.ID,
			UserID: task.UserID,
			CategoryID:  task.CategoryID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Priority:    task.Priority.String(),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	}, nil
}