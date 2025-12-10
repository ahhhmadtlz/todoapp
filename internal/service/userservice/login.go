package userservice

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {
	const op = richerror.Op("userservice.Login")

	user, err := s.repo.GetUserByPhoneNumber(ctx ,req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).WithMeta("phone_number", req.PhoneNumber)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("username or password is incorrect")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("un")
	}

	return param.LoginResponse{
		User: param.UserInfo{
		 ID:          user.ID,
		 PhoneNumber: user.PhoneNumber,
		 Name:        user.Name,
		},
		Tokens: param.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}