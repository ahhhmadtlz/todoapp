package categoryvalidator

import (
	"context"
	"todoapp/internal/entity"
)

type Repository interface {
GetCategoryByName(ctx context.Context, userID uint, name string) (entity.Category, error)
//                                    
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}