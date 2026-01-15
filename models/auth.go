package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	AccessLevel string `json:"acl"`
	UserID      string `json:"usr"`
	jwt.RegisteredClaims
}

const (
	ACCESS_TOKEN_COOKIE  string = "access_token"
	REFRESH_TOKEN_COOKIE string = "refresh_token"
	AUTH_HEADER          string = "Authorization"
)

var ErrJWTExpired error = errors.New("jwt expired")
var ErrJWTInvalid error = errors.New("invalid jwt")
