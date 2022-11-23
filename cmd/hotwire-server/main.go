package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	assets "github.com/wolfeidau/echo-esbuild-middleware"
	templates "github.com/wolfeidau/echo-go-templates"
	middleware "github.com/wolfeidau/echo-middleware"
	"github.com/wolfeidau/hotwire-golang-website/internal/app"
	"github.com/wolfeidau/hotwire-golang-website/internal/flags"
	"github.com/wolfeidau/hotwire-golang-website/internal/server"
	"github.com/wolfeidau/hotwire-golang-website/views"
)

var cfg = new(flags.ServerAPI)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

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

	e.Logger.SetOutput(io.Discard)

	e.Use(middleware.ZeroLogWithConfig(
		middleware.ZeroLogConfig{
			Fields: flds,
			Level:  zerolog.InfoLevel,
		},
	))

	e.Use(middleware.ZeroLogRequestLog())
	e.Use(echomiddleware.Gzip())

	hotwire := server.NewHotwire()

	server.RegisterHandlers(e, hotwire)

	// register the asset bundler which will build then serve any asset files
	e.Use(assets.BundlerWithConfig(assets.BundlerConfig{
		EntryPoints:     []string{"assets/src/index.ts"},
		Outfile:         "bundle.js",
		InlineSourcemap: cfg.Local,
		Define: map[string]string{
			"process.env.NODE_ENV": `"production"`,
		},
		OnBuild: func(result api.BuildResult, timeTaken time.Duration) {
			log.Info().Str("timeTaken", timeTaken.String()).Msg("bundle build complete")

			if len(result.Errors) > 0 {
				log.Fatal().Fields(map[string]interface{}{
					"errors": result.Errors,
				}).Msg("failed to build assets")
			}
		},
		OnRequest: func(req *http.Request, contentLength, code int, timeTaken time.Duration) {
			log.Ctx(req.Context()).Info().Str("path", req.URL.Path).Int("code", code).Str("timeTaken", timeTaken.String()).Msg("asset served")
		},
	}))

	log.Info().Str("addr", cfg.Addr).Str("cert", cfg.CertFile).Msg("listing")
	log.Error().Err(e.StartTLS(cfg.Addr, cfg.CertFile, cfg.KeyFile)).Msg("failed to bind port")
}
