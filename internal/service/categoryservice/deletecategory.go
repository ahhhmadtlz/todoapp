package categoryservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) DeleteCategory(ctx context.Context, ID uint ,UserID uint) (param.DeleteCategoryResponse ,error){
	const op=richerror.Op("categoryservice.DeleteCategory")

	_,err:=s.repo.GetCategoryByID(ctx,ID,UserID)

	if err !=nil{
		return param.DeleteCategoryResponse{Success: false},richerror.New(op).WithErr(err)
	}

	if err:=s.repo.DeleteCategory(ctx,ID,UserID);err!=nil{
		return param.DeleteCategoryResponse{Success: false},richerror.New(op).WithErr(err)
	}
	return  param.DeleteCategoryResponse{
		Success: true,
	},nil


}