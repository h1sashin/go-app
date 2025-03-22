package service

import (
	"github.com/h1sashin/go-app/config"
	"github.com/h1sashin/go-app/db"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	AuthService AuthService
	UserService UserService
	JWTService  JWTService
}

func NewService(db *pgx.Conn, queries *db.Queries, cfg *config.Config) *Service {
	jwtService := NewJWTService(cfg)
	userService := NewUserService(db, queries)
	authService := NewAuthService(db, queries, cfg, userService, jwtService)

	return &Service{
		AuthService: authService,
		UserService: userService,
		JWTService:  jwtService,
	}
}
