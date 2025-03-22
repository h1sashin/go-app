package db

import (
	"context"
	"fmt"

	"github.com/h1sashin/go-app/config"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func NewDB(cfg *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Connected to database successfully")
	return conn, nil
}
