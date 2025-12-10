package categoryservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) GetCategoryByID(ctx context.Context, ID uint, userID uint)(param.GetCategoryResponse,error){


	const op = richerror.Op("categoryservice.GetCategoryByID")
	
	cat, err := s.repo.GetCategoryByID(ctx, ID,userID)

	if err != nil {
		return param.GetCategoryResponse{}, richerror.New(op).WithErr(err)
	}
	
	return param.GetCategoryResponse{
		Category:param.CategoryInfo{
			ID:cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
		}}, nil
}