package htmx

import (
	"embed"
	"my-project-name/app/adapter/in/htmx/module"
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"

	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newIndex,
	server.NewServer,
	module.NewAppModule)

type index struct {
	server    server.Server
	appModule module.IModule
	URL       string
	PushURL   bool
	HTML      string
}

//go:embed *.html
var html embed.FS

//go:embed *.css
var css embed.FS

func newIndex(
	s server.Server,
	appModule module.IModule,
) (index, error) {
	view := index{
		appModule: appModule,
	}
	if err := s.TemplateRegistry(css, "index.css"); err != nil {
		return index{}, err
	}
	if err := s.TemplateRegistry(html, "index.html"); err != nil {
		return index{}, err
	}
	cssHandler := echo.WrapHandler(http.FileServer(http.FS(css)))
	s.Router().GET("/index.css", cssHandler)
	s.Router().GET("/", view.handle)
	return view, nil
}

func (state index) handle(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", state)
}
