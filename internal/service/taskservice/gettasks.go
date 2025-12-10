package taskservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) GetAllTasks(ctx context.Context, userID uint) (param.GetTasksResponse, error) {
	const op = richerror.Op("taskservice.GetAllTasks")

	tasks, err := s.repo.GetAllTasks(ctx, userID)
	if err != nil {
		return param.GetTasksResponse{}, richerror.New(op).WithErr(err)
	}

	taskInfos := make([]param.TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		taskInfos = append(taskInfos, param.TaskInfo{
			ID:          task.ID,
			UserID: task.UserID,
			CategoryID:  task.CategoryID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Priority:    task.Priority.String(),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	return param.GetTasksResponse{
		Tasks: taskInfos,
	}, nil
}
