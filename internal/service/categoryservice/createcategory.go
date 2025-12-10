package categoryservice

import (
	"context"
	"todoapp/internal/entity"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) CreateCategory(ctx context.Context,req param.CreateCategoryRequest)(param.CreatecategoryResponse,error){
 const op=richerror.Op("categoryservice.CreateCategory")

 var desc string
	if req.Description != nil {
		desc = *req.Description
	}
 
 category:=entity.Category{
	Name:req.Name,
	UserID: req.UserID,
	Description: desc,
 }

 createdCategory,err:=s.repo.CreateCategory(ctx,category)

 if err !=nil{
	return param.CreatecategoryResponse{},richerror.New(op).WithErr(err)
 }


 return param.CreatecategoryResponse{
	Category: param.CategoryInfo{
		ID:createdCategory.ID,
		Name:category.Name,
		Description: category.Description,
	},
 },nil

}