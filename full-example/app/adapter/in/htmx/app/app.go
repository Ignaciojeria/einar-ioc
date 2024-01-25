package app

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewApp,
	server.NewServer,
)

type App struct {
	server server.Server
	URL    string
	HTML   string
}

//go:embed *.html
var html embed.FS

func NewApp(
	server server.Server,
) (App, error) {
	view := App{
		URL:  "/app",
		HTML: "app.html",
	}

	if err := server.TemplateRegistry(html, view.HTML); err != nil {
		return App{}, err
	}

	//s.Router().GET(view.URL, view.handle, echo.WrapMiddleware(middleware.MiddleWare))
	return view, nil
}

func (state App) Render(c echo.Context) error {
	return c.Render(http.StatusOK, state.HTML, state)
}
