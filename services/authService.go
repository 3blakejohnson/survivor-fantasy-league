package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/3blakejohnson/survivor-fantasy-league/dao"
	"github.com/3blakejohnson/survivor-fantasy-league/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateJWT(userID string, accessLevel string, expiresAt *time.Time) (string, error)
	VerifyJWT(tokenString string) (*models.AccessClaims, error)
	ValidateLogin(username string, password string) (*models.User, error)
	HashPassword(password string) (string, error)
	CreateUser(models.User) (*models.User, error)
}

type authService struct {
	SecretKey string
	DM        dao.DAOManager
}

func NewAuthService(secretKey string, dm dao.DAOManager) AuthService {
	return &authService{
		SecretKey: secretKey,
		DM:        dm,
	}
}

func (a *authService) CreateJWT(userID string, accessLevel string, expiresAt *time.Time) (string, error) {
	if expiresAt == nil {
		timeout := time.Now().Add(time.Minute * 15)
		expiresAt = &timeout
	}
	claims := models.AccessClaims{
		AccessLevel: accessLevel,
		UserID:      userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(*expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.SecretKey))
}

func (a *authService) VerifyJWT(tokenString string) (*models.AccessClaims, error) {
	claims := &models.AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(a.SecretKey), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
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

func (a *authService) ValidateLogin(username string, password string) (*models.User, error) {
	user, err := a.DM.User().GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}

func (a *authService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (a *authService) CreateUser(user models.User) (*models.User, error) {
	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" || user.PasswordHash == "" {
		return nil, errors.New("username and password are required")
	}

	var err error
	user.PasswordHash, err = a.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if err := a.DM.User().Create(user); err != nil {
		return nil, err
	}

	return &user, nil
}
