package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger quick way to create a zerolog logger
func NewLogger(debug bool) zerolog.Logger {
	lvl := zerolog.InfoLevel
	if debug {
		lvl = zerolog.DebugLevel
	}

	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).Level((lvl)).With().Stack().Logger()
}
