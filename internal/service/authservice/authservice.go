package authservice

import (
	"fmt"
	"strings"
	"time"
	"todoapp/internal/entity"
	"todoapp/internal/pkg/richerror"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Operations
	OpCreateAccessToken  = richerror.Op("authservice.CreateAccessToken")
	OpCreateRefreshToken = richerror.Op("authservice.CreateRefreshToken")
	OpParseToken         = richerror.Op("authservice.ParseToken")
	OpCreateToken        = richerror.Op("authservice.createToken")
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject         string        `koanf:"access_subject"`
	RefreshSubject        string        `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}


func New (cfg Config) Service {
	return  Service{
		config:cfg,
	}
}

func (s Service ) CreateAccessToken(user entity.User)(string ,error){
 return s.createToken(user.ID,user.Role,s.config.AccessSubject,s.config.AccessExpirationTime)
}

func (s Service)CreateRefreshToken(user entity.User)(string,error){

 return s.createToken(user.ID,user.Role,s.config.RefreshSubject,s.config.RefreshExpirationTime)

}

func (s Service) ParseToken(bearerToken string)(*Claims,error){

	tokenStr := strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	if tokenStr == "" {
		return nil, richerror.New(OpParseToken).
			WithMessage("empty or invalid token format").
			WithKind(richerror.KindInvalid)
	}

	key := []byte(s.config.SignKey)

	token, err := jwt.ParseWithClaims(
			tokenStr,
			&Claims{},
			func(token *jwt.Token) (any, error) {
				// Verify signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, richerror.New(OpParseToken).
						WithMessage("invalid signing method").
						WithKind(richerror.KindInvalid).
						WithMeta("algorithm", token.Header["alg"])
				}
				return key, nil
			},
	)
	
	if err !=nil {
		return nil,err
	}
  
 
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("✅ userID: %v\n", claims.UserID)
    fmt.Printf("⏰ expires at: %v\n", claims.RegisteredClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
	}
}



func (s Service) createToken(
	userID uint,
	role entity.Role,
	subject string,
	expireDuration time.Duration,
) (string, error) {
	// Set claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Role:   role,
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}