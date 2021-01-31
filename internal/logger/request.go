package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// RequestLoggerConfig defines the config for the request logger middleware
type RequestLoggerConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// RequestLoggerWithConfig returns a request logger middleware with config.
func RequestLoggerWithConfig(config RequestLoggerConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}

			log.Ctx(c.Request().Context()).Info().Fields(map[string]interface{}{
				"path":   req.URL.Path,
				"method": req.Method,
				"dur":    time.Now().Sub(start).String(),
				"status": res.Status,
				"length": res.Size,
				"ip":     c.RealIP(),
			}).Msg("processed request")

			return err
		}
	}
}
