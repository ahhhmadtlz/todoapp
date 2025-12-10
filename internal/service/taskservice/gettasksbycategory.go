package taskservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) GetTasksByCategory(ctx context.Context, userID uint,categoryID uint)(param.GetTasksResponse,error){
	const op=richerror.Op("taskservice.GetTasksByCategoryID")
// because we dont have any validation we need to check in layer service this is an expection for this matter
	_,err:=s.categoryRepo.GetCategoryByID(ctx,categoryID,userID)

	if err !=nil{
		return  param.GetTasksResponse{},richerror.New(op).WithErr(err)
	}

	tasks,err:=s.repo.GetTasksByCategory(ctx,userID,categoryID)

	if err!=nil{
		return  param.GetTasksResponse{},richerror.New(op).WithErr(err)
	}

	taskInfos:=make([]param.TaskInfo,0,len(tasks))

	for _,task:=range tasks{
		taskInfos=append(taskInfos, param.TaskInfo{
			ID:          task.ID,
			UserID: task.UserID,
			CategoryID:  task.CategoryID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Priority:    task.Priority.String(),
			Status: task.Status.String(),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

return param.GetTasksResponse{
		Tasks: taskInfos,
	}, nil
}

