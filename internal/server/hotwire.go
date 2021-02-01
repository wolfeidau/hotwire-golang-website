package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	turboStreamMedia = "text/vnd.turbo-stream.html"
)

var (
	upgrader = websocket.Upgrader{}
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

// Memory uses websockets
func (hw *Hotwire) Memory(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	ticker := time.NewTicker(3 * time.Second)
	mstats := new(runtime.MemStats)

	for range ticker.C {
		buf := new(bytes.Buffer)

		runtime.ReadMemStats(mstats)

		err := c.Echo().Renderer.Render(buf, "memory.turbo-stream.html", mstats, c)
		if err != nil {
			log.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to build message")
			break
		}

		err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			log.Ctx(c.Request().Context()).Error().Err(err).Msg("send failed")
			break
		}
	}

	return nil
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

	seq := 0

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

			err = writeMessage(res.Writer, seq, "message", buf.String())
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			flusher.Flush()

			seq++
		}
	}
}

// writeMessage this constructs an SSE compatible message with a sequence, and
// line breaks from the output of a template
//
// This looks something like this:
//
//   event: message
//   id: 6
//   data: <turbo-stream action="replace" target="load">
//   data:     <template>
//   data:         <span id="load">04:20:13: 1.9</span>
//   data:     </template>
//   data: </turbo-stream>
//
func writeMessage(w io.Writer, id int, event, message string) error {

	_, err := fmt.Fprintf(w, "event: %s\nid: %d\n", event, id)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewBufferString(message))
	for scanner.Scan() {
		_, err = fmt.Fprintf(w, "data: %s\n", scanner.Text())
		if err != nil {
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\n")
	if err != nil {
		return err
	}

	return nil
}
