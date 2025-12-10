package userservice

import (
	"context"
	"todoapp/internal/entity"
)

type Repository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, error)
	GetUserByPhoneNumber(ctx context.Context, phone string) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthGenerator
}

func New(authGenerator AuthGenerator,repo Repository,) Service {
	return Service{
		auth:authGenerator,
		repo: repo,
	}
}
