package taskservice

import (
	"context"
	"todoapp/internal/pkg/richerror"
)

func (s Service) DeleteTask(ctx context.Context,id uint,userID uint) error {
	const op=richerror.Op("taskservice.DeleteTask")
	err:=s.repo.DeleteTask(ctx,id,userID)

	if err !=nil{
		return  richerror.New(op).WithErr(err)
	}
	return nil
}