package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/h1sashin/go-app/config"
	"github.com/h1sashin/go-app/db"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	DB          *pgx.Conn
	queries     *db.Queries
	Cfg         *config.Config
	UserService UserService
	JWTService  JWTService
}

type AuthService interface {
	SignIn(ctx context.Context, email, password string) (*db.User, *Tokens, error)
	SignUp(ctx context.Context, email, password string, role db.Role) (*db.User, *Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error)
}

func NewAuthService(db *pgx.Conn, queries *db.Queries, cfg *config.Config, userService UserService, jwtService JWTService) AuthService {
	return &authService{
		DB:          db,
		queries:     queries,
		Cfg:         cfg,
		UserService: userService,
		JWTService:  jwtService,
	}
}

func (s *authService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.Cfg.SaltOrRounds)
	if err != nil {
		return "", ErrInternal
	}
	return string(hashedPassword), nil
}

func (s *authService) comparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}

func (s *authService) generateTokens(userID uuid.UUID, email string, role db.Role) (*Tokens, error) {
	accessToken, err := s.JWTService.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*db.User, *Tokens, error) {
	user, err := s.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrNotFound
	}

	if err := s.comparePassword(user.Password, password); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, ErrInternal
	}

	return user, tokens, nil
}

func (s *authService) SignUp(ctx context.Context, email, password string, role db.Role) (*db.User, *Tokens, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, nil, ErrInternal
	}

	defer tx.Rollback(ctx)

	ctxWithTx := db.InjectTx(ctx, tx)

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.UserService.CreateUser(ctxWithTx, email, hashedPassword, db.Role(role))
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, ErrInternal
	}

	return user, tokens, nil

}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error) {
	claims, err := s.JWTService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	tokens, err := s.generateTokens(claims.UserID, claims.Email, db.Role(claims.Role))
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
