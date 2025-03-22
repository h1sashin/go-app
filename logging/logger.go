package logging

import (
	"os"
	"time"

	"github.com/h1sashin/go-app/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(logLevel config.LogLevel) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	level := zerolog.InfoLevel

	switch logLevel {
	case config.LogLevel(config.Debug):
		level = zerolog.DebugLevel
	case config.LogLevel(config.Warn):
		level = zerolog.WarnLevel
	case config.LogLevel(config.Error):
		level = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	log.Info().Msgf("Logger initialized with level: %s", level.String())
}

func Logger(component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}
