package htmx

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/heimdalr/dag"

	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	NewIndex,
	server.NewServer)

type Index struct {
	URL string
	Dag *dag.DAG
	//PushURL       bool
	//DefaultModule string
	HTML string
	CSS  string
}

//go:embed *.html
var html embed.FS

//go:embed *.css
var css embed.FS

func NewIndex(s server.Server) (Index, error) {
	view := Index{
		URL: "/",
		//DefaultModule: "/app",
		//PushURL:       false,
		HTML: "index.html",
		CSS:  "index.css",
	}
	if err := s.TemplateRegistry(css, view.CSS); err != nil {
		return Index{}, err
	}
	if err := s.TemplateRegistry(html, view.HTML); err != nil {
		return Index{}, err
	}
	//cssHandler := echo.WrapHandler(http.FileServer(http.FS(css)))
	//s.Router().GET(view.URL+view.CSS, cssHandler)
	//s.Router().GET(view.URL, view.handle)
	return view, nil
}

func (state Index) Render(c echo.Context) error {
	err := c.Render(http.StatusOK, state.HTML, state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}
