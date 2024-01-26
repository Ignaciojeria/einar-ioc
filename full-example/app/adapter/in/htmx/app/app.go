package app

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"my-project-name/app/infrastructure/uirouter"
	"net/http"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewApp,
	server.NewServer,
)

type App struct {
	server      server.Server
	ActiveRoute uirouter.Route
	URL         string
	HTML        string
	Target      string
}

//go:embed *.html
var html embed.FS

func NewApp(
	server server.Server,
) (App, error) {
	view := App{
		URL:    "/app",
		HTML:   "app.html",
		Target: strings.ReplaceAll(uuid.NewString(), "-", ""),
	}
	if err := server.TemplateRegistry(html, view.HTML); err != nil {
		return App{}, err
	}
	return view, nil
}

func (state App) Render(c echo.Context) error {
	state.ActiveRoute = c.Get("activeRoute").(uirouter.Route)
	return c.Render(http.StatusOK, state.HTML, state)
}
