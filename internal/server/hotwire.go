package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	turboStreamMedia = "text/vnd.turbo-stream.html"
)

// Hotwire hotwire handlers which demonstate some of the capabilities
type Hotwire struct{}

// NewHotwire new hotwire handlers
func NewHotwire() *Hotwire {
	return &Hotwire{}
}

// Index using a template build the index page
func (hw *Hotwire) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// Greeting process greetings using a query parameter
func (hw *Hotwire) Greeting(c echo.Context) error {
	name := c.QueryParam("person")

	return c.Render(http.StatusOK, "greeting.html", map[string]interface{}{
		"person": name,
	})
}

// Pinger process form POSTs which are either returned as a partial or full page dependending on the media type
func (hw *Hotwire) Pinger(c echo.Context) error {
	req := c.Request()

	ctype := req.Header.Get(echo.HeaderContentType)
	accept := req.Header.Get(echo.HeaderAccept)

	log.Ctx(c.Request().Context()).Debug().Str("contentType", ctype).Str("accept", accept).Msg("pinger")

	if strings.HasPrefix(accept, turboStreamMedia) {

		c.Response().Header().Set(echo.HeaderContentType, turboStreamMedia)

		return c.Render(http.StatusOK, "ping.turbo-stream.html", map[string]interface{}{
			"pingTime": 0,
		})
	}

	return c.Render(http.StatusOK, "ping.html", map[string]interface{}{
		"pingTime": 0,
	})
}

// Load use http2 server sent events to stream load information
func (hw *Hotwire) Load(c echo.Context) error {
	req := c.Request()
	res := c.Response()

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported!")
	}

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-req.Context().Done():
			return c.NoContent(http.StatusOK)
		case t := <-ticker.C:

			buf := new(bytes.Buffer)

			err := c.Echo().Renderer.Render(buf, "load.turbo-stream.html", map[string]interface{}{
				"at":  t.Format("15:04:05"),
				"avg": runtime.NumGoroutine(),
			}, c)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			log.Ctx(req.Context()).Debug().Str("buf", buf.String()).Msg("event")

			_, err = writeMessageWithoutNewline(res.Writer, buf.String())
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			flusher.Flush()
		}
	}
}

// This method strips line feeds from the templated result to enable single line responses on the events stream
func writeMessageWithoutNewline(w io.Writer, message string) (int, error) {
	return fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(message, "\n", ""))
}
