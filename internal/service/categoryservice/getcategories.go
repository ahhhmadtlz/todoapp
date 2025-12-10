package categoryservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"
)

func (s Service) GetAllCategories(ctx context.Context, userID uint) (param.GetCategoriesResponse, error) {
	const op = richerror.Op("categoryservice.GetAllCategories")
	
	categories, err := s.repo.GetAllCategories(ctx, userID)
	if err != nil {
		return param.GetCategoriesResponse{}, richerror.New(op).WithErr(err)
	}
	
	categoryInfos := make([]param.CategoryInfo, 0, len(categories))
	for _, cat := range categories {
		categoryInfos = append(categoryInfos, param.CategoryInfo{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
		})
	}
	
	return param.GetCategoriesResponse{
		Categories: categoryInfos,
	}, nil
}