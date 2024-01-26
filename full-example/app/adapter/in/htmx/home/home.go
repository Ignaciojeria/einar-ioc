package home

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewHome,
	server.NewServer,
)

type Home struct {
	server server.Server
	URL    string
	HTML   string
}

//go:embed *.html
var html embed.FS

func NewHome(
	server server.Server,
) (Home, error) {
	view := Home{
		URL:  "/home",
		HTML: "home.html",
	}
	if err := server.TemplateRegistry(html, view.HTML); err != nil {
		return Home{}, err
	}
	return view, nil
}

func (state Home) Render(c echo.Context) error {
	return c.Render(http.StatusOK, state.HTML, state)
}
