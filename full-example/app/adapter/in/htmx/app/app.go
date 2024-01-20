package app

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/donseba/go-htmx/middleware"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewAppModule,
	server.NewServer,
)

type AppModule struct {
	server  server.Server
	URL     string
	PushURL bool
	HTML    string
}

//go:embed *.html
var html embed.FS

func NewAppModule(s server.Server) (AppModule, error) {
	view := AppModule{
		URL:     "/app",
		PushURL: false,
		HTML:    "app.html",
		server:  s,
	}
	if err := s.TemplateRegistry(html, view.HTML); err != nil {
		return AppModule{}, err
	}

	s.Router().GET(view.URL, view.handle, echo.WrapMiddleware(middleware.MiddleWare))
	return view, nil
}

func (state AppModule) handle(c echo.Context) error {
	h := state.
		server.HTMX().
		NewHandler(c.Response().Writer, c.Request())

	// check if the request is a htmx request
	if h.IsHxRequest() {
		// do something
	}

	// check if the request is boosted
	if h.IsHxBoosted() {
		// do something
	}

	// check if the request is a history restore request
	if h.IsHxHistoryRestoreRequest() {
		// do something
	}

	// check if the request is a prompt request
	if h.RenderPartial() {
		// do something
	}

	// set the headers for the response, see docs for more options
	//h.PushURL("http://push.url")
	//h.ReTarget("#ReTarged")
	return c.Render(http.StatusOK, state.HTML, state)
}
