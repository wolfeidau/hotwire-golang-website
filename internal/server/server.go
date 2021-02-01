package server

import "github.com/labstack/echo/v4"

// EchoRouter This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// ServerInterface implemented by the handlers
type ServerInterface interface {
	Index(echo.Context) error
	Greeting(echo.Context) error
	Pinger(echo.Context) error
	Memory(echo.Context) error
	Load(echo.Context) error
}

// RegisterHandlers register the handlers
func RegisterHandlers(router EchoRouter, si ServerInterface, m ...echo.MiddlewareFunc) {
	router.GET("/", si.Index).Name = "index"
	router.GET("/greeting", si.Greeting).Name = "greeting"
	router.POST("/pinger", si.Pinger).Name = "pinger"
	router.GET("/memory", si.Memory).Name = "memory"
	router.GET("/load", si.Load).Name = "load"
}
