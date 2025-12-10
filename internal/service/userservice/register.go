package userservice

import (
	"context"
	"todoapp/internal/entity"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)



func (s Service) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {
	const op = richerror.Op("userService.Register")



	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	if err != nil {
		return param.RegisterResponse{}, richerror.New(op).
			WithMessage("failed to hash password").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	user := entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		Role:        entity.UserRole, // or entity.RoleUser
	}

	// Save user
	createdUser, err := s.repo.RegisterUser(ctx, user)
	if err != nil {
		// Wrap repository errors in richerror
		return param.RegisterResponse{}, richerror.New(op).
			WithMessage("failed to register user").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	return param.RegisterResponse{
		User: param.UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}
