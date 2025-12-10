package uservalidator

import (
	"context"
	"regexp"
	"todoapp/internal/entity"
)


var (
	reUpper         = regexp.MustCompile(`[A-Z]`)
	reLower         = regexp.MustCompile(`[a-z]`)
	reDigit         = regexp.MustCompile(`[0-9]`)
	reSpecial       = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|` + "`" + `]`)
	phoneNumberRegex = regexp.MustCompile("^09[0-9]{9}$")
	passwordRegex    = regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`)
)

type Repository interface {
	IsPhoneNumberUnique(ctx context.Context, phoneNumber string) (bool, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}