package taskservice

import (
	"context"
	"todoapp/internal/entity"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) UpdateTask(ctx context.Context,req param.UpdateTaskRequest)(param.UpdateTaskResponse,error){
	const op =richerror.Op("taskservice.UpdateTask")

	existingTask,err:=s.repo.GetTaskByID(ctx,req.ID,req.UserID)

	if err !=nil{
		return param.UpdateTaskResponse{},richerror.New(op).WithErr(err)
	}

	if req.Title != nil && *req.Title != "" {
   existingTask.Title = *req.Title
  	}
  if req.Description != nil {
      existingTask.Description = *req.Description
   }
  if req.CategoryID != nil {
      existingTask.CategoryID = *req.CategoryID
   }
  if req.Priority != nil && *req.Priority != "" {
		existingTask.Priority = entity.MapToPriorityEntity(*req.Priority)
	} 
   if req.Status != nil && *req.Status != "" {
       existingTask.Status = entity.MapToStatusEntity(*req.Status)
   }
   if req.DueDate != nil {
        existingTask.DueDate = req.DueDate
   }

	updatedTask,err:=s.repo.UpdateTask(ctx,existingTask)

	if err!=nil{
		return  param.UpdateTaskResponse{},richerror.New(op).WithErr(err)
	}


	
	return param.UpdateTaskResponse{
		Task: param.TaskInfo{
          ID:          updatedTask.ID,
					UserID: updatedTask.UserID,
          CategoryID:  updatedTask.CategoryID,
          Title:       updatedTask.Title,
          Description: updatedTask.Description,
          DueDate:     updatedTask.DueDate,
          Priority:    updatedTask.Priority.String(),
          Status:      updatedTask.Status.String(),
			 CreatedAt:   updatedTask.CreatedAt,
			 UpdatedAt:   updatedTask.UpdatedAt,
        },
	},nil

}