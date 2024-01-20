package server

import (
	"embed"
	"fmt"
	"io"
	"my-project-name/app/configuration"
	"text/template"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewServer, configuration.NewConf)

type Server struct {
	c configuration.Conf
	e *echo.Echo
	t *template.Template
	h *htmx.HTMX
}

func NewServer(c configuration.Conf) Server {
	e := echo.New()
	return Server{
		e: e,
		c: c,
		t: template.New(""),
		h: htmx.New(),
	}
}

func (s Server) Start() {
	s.printRoutes()
	s.e.Renderer = &templateRegistry{
		templates: s.t,
	}
	s.e.Start(":" + s.c.Port)
}

func (s Server) printRoutes() {
	for _, route := range s.e.Routes() {
		fmt.Printf("Method: %v, Path: %v, Name: %v\n", route.Method, route.Path, route.Name)
	}
}

func (s Server) Router() *echo.Echo {
	return s.e
}

func (s Server) HTMX() *htmx.HTMX {
	return s.h
}

func (s Server) TemplateRegistry(fs embed.FS, pattern string) error {
	t, err := s.t.ParseFS(fs, pattern)
	s.t = t
	return err
}

type templateRegistry struct {
	templates *template.Template
}

func (t *templateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
