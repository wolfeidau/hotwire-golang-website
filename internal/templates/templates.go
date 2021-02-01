package templates

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var templateFuncs = template.FuncMap{
	"getTime": func() string {
		return time.Now().Format("15:04:05")
	},
}

// Template stores the meta data for each template, and whether it uses a layout
type Template struct {
	layout   string
	name     string
	template *template.Template
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates map[string]*Template
}

// New setup a new template renderer
func New() *TemplateRenderer {
	return &TemplateRenderer{
		templates: make(map[string]*Template),
	}
}

// AddWithLayout register one or more templates using the provided layout
func (t *TemplateRenderer) AddWithLayout(fsys fs.FS, layout string, patterns ...string) error {
	filenames, err := readFileNames(fsys, patterns...)
	if err != nil {
		return errors.Wrap(err, "failed to list using file pattern")
	}

	for _, f := range filenames {

		tname := path.Base(f)
		lname := path.Base(layout)

		log.Debug().Str("filename", tname).Str("layout", layout).Msg("register template")
		t.templates[tname] = &Template{
			layout:   lname,
			name:     tname,
			template: template.Must(template.New(tname).Funcs(templateFuncs).ParseFS(fsys, layout, f)),
		}
	}

	return nil
}

// Add add a template to the registry
func (t *TemplateRenderer) Add(fsys fs.FS, patterns ...string) error {
	filenames, err := readFileNames(fsys, patterns...)
	if err != nil {
		return errors.Wrap(err, "failed to read file names using file pattern")
	}

	for _, f := range filenames {
		tname := path.Base(f)

		log.Debug().Str("filename", tname).Msg("register message")
		t.templates[tname] = &Template{
			name:     tname,
			template: template.Must(template.New(tname).Funcs(templateFuncs).ParseFS(fsys, f)),
		}
	}

	return nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	log.Ctx(c.Request().Context()).Debug().Str("name", name).Msg("Render")

	tmpl, ok := t.templates[name]
	if !ok {
		log.Ctx(c.Request().Context()).Error().Str("name", name).Msg("template not found")

		return c.NoContent(http.StatusInternalServerError)
	}

	// use the name of the template, or layout if it exists
	execName := tmpl.name
	if tmpl.layout != "" {
		execName = tmpl.layout
	}

	start := time.Now()
	err := tmpl.template.ExecuteTemplate(w, execName, data)
	if err != nil {
		log.Ctx(c.Request().Context()).Error().Err(err).Str("name", tmpl.name).Str("layout", tmpl.layout).Msg("render template failed")
		return err
	}

	log.Ctx(c.Request().Context()).Debug().Str("name", tmpl.name).Str("dur", time.Since(start).String()).Str("layout", tmpl.layout).Msg("execute template")

	return nil
}

func readFileNames(fsys fs.FS, patterns ...string) ([]string, error) {
	var filenames []string

	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list using file pattern")
		}

		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		filenames = append(filenames, list...)
	}

	return filenames, nil
}
