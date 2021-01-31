package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	middleware "github.com/wolfeidau/echo-middleware"
	"github.com/wolfeidau/hotwire-golang-website/internal/app"
	"github.com/wolfeidau/hotwire-golang-website/internal/flags"
	"github.com/wolfeidau/hotwire-golang-website/internal/logger"
	"github.com/wolfeidau/hotwire-golang-website/internal/server"
	"github.com/wolfeidau/hotwire-golang-website/internal/templates"
	"github.com/wolfeidau/hotwire-golang-website/views"
)

var cfg = new(flags.ServerAPI)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	var output io.Writer = os.Stderr

	// enable pretty output for local development
	if cfg.Local {
		log.Logger = logger.NewLogger().Level(cfg.ZerologLevel())
		output = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}
	}

	flds := map[string]interface{}{"commit": app.Commit, "buildDate": app.BuildDate}

	e := echo.New()

	render := templates.New()

	err := render.AddWithLayout(views.Content, "layouts/base.html", "templates/*.html")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load render")
	}

	err = render.Add(views.Content, "messages/*.html")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load render")
	}

	e.Renderer = render

	e.Logger.SetOutput(ioutil.Discard)

	e.Use(middleware.ZeroLogWithConfig(
		middleware.ZeroLogConfig{
			Fields: flds,
			Output: output,
			Level:  cfg.ZerologLevel(),
		},
	))

	e.Use(middleware.ZeroLogRequestLog())
	e.Use(echomiddleware.Gzip())

	hotwire := server.NewHotwire()

	server.RegisterHandlers(e, hotwire)

	log.Info().Str("port", cfg.Port).Str("cert", cfg.CertFile).Msg("listing")
	log.Error().Err(e.StartTLS(fmt.Sprintf(":%s", cfg.Port), cfg.CertFile, cfg.KeyFile)).Msg("failed to bind port")
}
