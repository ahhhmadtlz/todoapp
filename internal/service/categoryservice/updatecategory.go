package categoryservice

import (
	"context"
	"strings"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) UpdateCategory(ctx context.Context,req param.UpdateCategoryRequest)(param.UpdateCategoryResponse,error){
	const op=richerror.Op("categoryservice.UpdateCategory")

	existingCategory,err:=s.repo.GetCategoryByID(ctx,req.ID,req.UserID)

	if err !=nil{
		return param.UpdateCategoryResponse{},richerror.New(op).WithErr(err)
	}

	if req.Name != nil && *req.Name != "" {
		name := strings.TrimSpace(*req.Name)
		existingCategory.Name = name
	}

	if req.Description !=nil{
		description:=strings.TrimSpace(*req.Description)
		existingCategory.Description=description
	}

	updatedCategory,err:=s.repo.UpdateCategory(ctx,existingCategory)

	if err!=nil{
		return  param.UpdateCategoryResponse{},richerror.New(op).WithErr(err)
	}

	return  param.UpdateCategoryResponse{
		Category:param.CategoryInfo{
			ID:updatedCategory.ID,
			Name:updatedCategory.Name,
			Description: updatedCategory.Description,
		},
	},nil
	
}