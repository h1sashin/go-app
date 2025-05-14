package logging

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/h1sashin/go-app/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(cfg *config.Config) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	level := zerolog.InfoLevel

	switch cfg.LogLevel {
	case config.Debug:
		level = zerolog.DebugLevel
	case config.Warn:
		level = zerolog.WarnLevel
	case config.Error:
		level = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(level)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDsn,
		Environment: string(cfg.AppEnv),
	}); err != nil {
		log.Error().Err(err).Msg("Failed to initialize Sentry")
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	log.Info().Msgf("Logger initialized with level: %s", level.String())
}

func Logger(component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}
