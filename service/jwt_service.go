package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/h1sashin/go-app/config"
	"github.com/h1sashin/go-app/db"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, email string, role db.Role) (string, error)
	ValidateAccessToken(token string) (*Claims, error)
	GenerateRefreshToken(userID uuid.UUID, email string, role db.Role) (string, error)
	ValidateRefreshToken(token string) (*Claims, error)
}

type jwtService struct {
	AccessSecretKey  []byte
	AccessDuration   time.Duration
	RefreshSecretKey []byte
	RefreshDuration  time.Duration
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		AccessSecretKey:  []byte(cfg.AccessSecretKey),
		AccessDuration:   cfg.AccessDuration,
		RefreshSecretKey: []byte(cfg.RefreshSecretKey),
		RefreshDuration:  cfg.AccessDuration,
	}
}

func (s *jwtService) GenerateAccessToken(userID uuid.UUID, email string, role db.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.AccessSecretKey)
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInternal
		}
		return s.AccessSecretKey, nil
	})

	if err != nil {
		return nil, ErrInternal
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrUnauthorized
	}

	return claims, nil
}

func (s *jwtService) GenerateRefreshToken(userID uuid.UUID, email string, role db.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.RefreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.RefreshSecretKey)
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.RefreshSecretKey, nil
	})

	if err != nil {
		return nil, ErrInternal
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrUnauthorized
	}

	return claims, nil
}
