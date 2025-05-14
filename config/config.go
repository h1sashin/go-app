package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type LogLevel string

const (
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

type Config struct {
	AppEnv           Environment
	AppPort          int
	DatabaseUrl      string
	SentryDsn        string
	LogLevel         LogLevel
	AccessSecretKey  string
	AccessDuration   time.Duration
	RefreshSecretKey string
	RefreshDuration  time.Duration
	SaltOrRounds     int
}

func Load() (*Config, error) {
	err := godotenv.Load(".env.local")

	if err != nil {
		panic(err)
	}

	appPort, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	accessDuration, _ := time.ParseDuration(getEnv("JWT_ACCESS_DURATION", "5m"))
	refreshDuration, _ := time.ParseDuration(getEnv("JWT_REFRESH_DURATION", "7d"))
	saltOrRounds, _ := strconv.ParseInt(getEnv("SALT_OR_ROUNDS", "2"), 10, 64)

	return &Config{
		AppEnv:           Environment(getEnv("APP_ENV", "development")),
		AppPort:          appPort,
		DatabaseUrl:      getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"),
		SentryDsn:        getEnv("SENTRY_DSN", ""),
		LogLevel:         LogLevel(getEnv("LOG_LEVEL", "info")),
		AccessSecretKey:  getEnv("JWT_ACCESS_SECRET", "access_secret"),
		AccessDuration:   accessDuration,
		RefreshSecretKey: getEnv("JWT_REFRESH_SECRET", "refresh_secret"),
		RefreshDuration:  refreshDuration,
		SaltOrRounds:     int(saltOrRounds),
	}, nil

}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
