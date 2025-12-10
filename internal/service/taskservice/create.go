package taskservice

import (
	"context"
	"todoapp/internal/entity"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) CreateTask(ctx context.Context, req param.CreateTaskRequest)(param.CreateTaskResponse,error){
	const op=richerror.Op("taskservice.CreateTask")

	var desc string

	if req.Description !=nil{
		desc =*req.Description
	}

	task:=entity.Task{
		UserID: req.UserID,
		CategoryID:req.CategoryID,
		Title: req.Title,
		Description: desc,
		DueDate: req.DueDate,
		Priority: entity.MapToPriorityEntity(req.Priority),
			Status: entity.MapToStatusEntity(req.Status),
	}

	createdTask,err:=s.repo.CreateTask(ctx,task)

	if err !=nil{
		return param.CreateTaskResponse{},richerror.New(op).WithErr(err)
	}

	return  param.CreateTaskResponse{
		Task: param.TaskInfo{
			ID:createdTask.ID,
			UserID: createdTask.UserID,
			CategoryID: createdTask.CategoryID,
			Title: createdTask.Title,
			Description: createdTask.Description,
			DueDate: createdTask.DueDate,
			Priority: createdTask.Priority.String(),
			Status: createdTask.Status.String(),
			CreatedAt:   createdTask.CreatedAt,
			UpdatedAt:   createdTask.UpdatedAt,
		},
	},nil

}