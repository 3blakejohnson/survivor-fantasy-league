package services

import (
	"errors"
	"time"

	"github.com/3blakejohnson/survivor-fantasy-league/models"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	VerifyJWT(tokenString string) (*models.AccessClaims, error)
}

type authService struct {
	SecretKey string
}

func NewAuthService(secretKey string) AuthService {
	return &authService{
		SecretKey: secretKey,
	}
}

func (a *authService) VerifyJWT(tokenString string) (*models.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		models.AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.SecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, models.ErrJWTInvalid
	}

	if claims, ok := token.Claims.(*models.AccessClaims); ok {
		now := time.Now()
		exp := claims.ExpiresAt.Time
		if exp.Before(now) {
			return nil, models.ErrJWTExpired
		}
		return claims, nil
	}
	return nil, errors.New("missing access claims")
}
