package taskservice

import (
	"context"
	"todoapp/internal/entity"
)

type Repository interface {
	CreateTask(ctx context.Context,task entity.Task)(entity.Task,error)
	GetTaskByID(ctx context.Context,id uint,userID uint)(entity.Task,error)
	GetAllTasks(ctx context.Context,userID uint)([]entity.Task,error)
	GetTasksByCategory(ctx context.Context,userID uint,categoryID uint) ([]entity.Task, error)
	UpdateTask(ctx context.Context,task entity.Task)(entity.Task, error)
	DeleteTask(ctx context.Context,id uint , userID uint)error
}

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, categoryID uint, userID uint) (entity.Category, error)
}

type Service struct {
	repo Repository
	categoryRepo CategoryRepository
}


func New(repo Repository,categoryRepo CategoryRepository)Service{
	return  Service{
		repo:repo,
		categoryRepo: categoryRepo,
	}
}