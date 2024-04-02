package app

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"my-project-name/app/infrastructure/uicomponent"
	"my-project-name/app/infrastructure/uirouter"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewApp,
	server.NewServer,
)

type App struct {
	uicomponent.Component
}

//go:embed *.html
var html embed.FS

func NewApp(
	server server.Server,
) (App, error) {
	view := App{
		Component: uicomponent.Component{
			URL:    "/app",
			HTML:   "app.html",
			Target: uirouter.NewSelectorTarget(),
		},
	}
	if err := server.TemplateRegistry(html, view.HTML); err != nil {
		return App{}, err
	}
	return view, nil
}

func (state App) Render(c echo.Context) error {
	return c.Render(
		http.StatusOK,
		state.HTML,
		state.WithContext(c))
}
