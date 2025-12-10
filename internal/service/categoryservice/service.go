package categoryservice

import (
	"context"
	"todoapp/internal/entity"
)

type Repository interface {
	CreateCategory(ctx context.Context,category entity.Category)(entity.Category,error)
	GetCategoryByID(ctx context.Context,id uint,userID uint)(entity.Category,error)
	GetAllCategories(ctx context.Context,userID uint)([]entity.Category,error)
	UpdateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, id uint,userID uint) error
}


type Service struct {
	repo Repository
}


func New(repo Repository) Service{
	return Service{
		repo:repo,
	}
}