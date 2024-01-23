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
	s server.Server) (App, error) {
	view := App{
		URL: "/app",
		//PushURL: false,
		HTML:   "app.html",
		server: s,
	}
	if err := s.TemplateRegistry(html, view.HTML); err != nil {
		return App{}, err
	}

	//s.Router().GET(view.URL, view.handle, echo.WrapMiddleware(middleware.MiddleWare))
	return view, nil
}

func (state App) Render(c echo.Context) error {
	h := state.
		server.HTMX().
		NewHandler(c.Response().Writer, c.Request())

	contextGraph := state.server.ContextGraph(c)

	if !h.IsHxRequest() {
		//TODO :  Enviar en la redirección el contexto del nodo desde dónde se redireccionó
		//TODO :  Ver la manera de representar multiples grafos para tener diferentes modulos de ruteo
		return c.Redirect(304, contextGraph.OrderedPaths[0])
		// do something
	}

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
	return c.Render(http.StatusOK, state.HTML, contextGraph)
}
