package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/h1sashin/go-app/db"
	"github.com/jackc/pgx/v5"
)

type userService struct {
	Queries *db.Queries
	DB      *pgx.Conn
}

type UserService interface {
	CreateUser(ctx context.Context, email, password string, role db.Role) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetUserByID(ctx context.Context, id string) (*db.User, error)
}

func NewUserService(db *pgx.Conn, queries *db.Queries) UserService {
	return &userService{
		DB:      db,
		Queries: queries,
	}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	tx, _ := db.ExtractTx(ctx)
	q := s.Queries
	if tx != nil {
		q = s.Queries.WithTx(tx)
	}

	user, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, ErrInternal
	}

	return &user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	tx, ok := db.ExtractTx(ctx)
	q := s.Queries
	if !ok {
		q = s.Queries.WithTx(tx)
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInternal
	}

	user, err := q.GetUserByID(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, ErrInternal
	}

	return &user, nil
}

func (s *userService) CreateUser(ctx context.Context, email, password string, role db.Role) (*db.User, error) {
	var err error
	tx, ok := db.ExtractTx(ctx)
	if !ok {
		tx, err = s.DB.Begin(ctx)
		if err != nil {
			return nil, ErrInternal
		}
		defer tx.Rollback(ctx)
	}

	q := s.Queries.WithTx(tx)

	userInDB, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInternal
	}
	if userInDB != nil {
		return nil, ErrConflict
	}

	user, err := q.CreateUser(ctx, db.CreateUserParams{Email: email, Password: password, Role: role})
	if err != nil {
		return nil, ErrInternal
	}

	if !ok {
		if err := tx.Commit(ctx); err != nil {
			return nil, ErrInternal
		}
	}

	return &user, nil
}
