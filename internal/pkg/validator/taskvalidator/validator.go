package taskvalidator

import (
	"context"
	"todoapp/internal/entity"
)

type Repository interface {

}

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context,categoryID uint,userID uint)(entity.Category,error)
}

type Validator struct {
	repo Repository
	categoryRepo CategoryRepository
}

func New(repo Repository,categoryRepo CategoryRepository) Validator {
	return Validator{repo: repo,categoryRepo: categoryRepo}
}
