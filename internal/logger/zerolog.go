package logger

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).With().Stack().Logger()
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	zlog := NewLogger()

	return zlog.WithContext(ctx)
}

func Middleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// configure zerolog in the context
		ctx := NewLoggerWithContext(c.Request().Context())

		// update the request
		req := c.Request().WithContext(ctx)

		// assign the new request to context
		c.SetRequest(req)

		return handlerFunc(c)
	}
}
