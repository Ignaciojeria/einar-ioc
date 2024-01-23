package server

import (
	"embed"
	"fmt"
	"io"
	"my-project-name/app/configuration"
	"text/template"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/donseba/go-htmx"
	"github.com/heimdalr/dag"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewServer, configuration.NewConf)

type Server struct {
	c            configuration.Conf
	e            *echo.Echo
	t            *template.Template
	h            *htmx.HTMX
	appDag       *dag.DAG
	contextGraph contextGraph
}

func NewServer(c configuration.Conf) Server {
	e := echo.New()
	return Server{
		e:      e,
		c:      c,
		t:      template.New(""),
		h:      htmx.New(),
		appDag: dag.NewDAG(),
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

func (s Server) AppDag() *dag.DAG {
	return s.appDag
}

func (s Server) ApiPrefix() string {
	return s.c.ApiPrefix
}

type templateRegistry struct {
	templates *template.Template
}

func (t *templateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type contextGraph struct {
	ctx          echo.Context
	OrderedPaths []string
	VertexPath   string
	EdgePath     string
}

func (v *contextGraph) Visit(vertex dag.Vertexer) {
	_, path := vertex.Vertex()
	v.appendOrderedPath(v.ctx, path.(string))
}

func (s Server) ContextGraph(ctx echo.Context) contextGraph {
	contextGraph := contextGraph{
		ctx: ctx,
	}
	s.appDag.OrderedWalk(&contextGraph)
	contextGraph.processCurrentPath(ctx)
	return contextGraph
}

func (v *contextGraph) appendOrderedPath(ctx echo.Context, path string) {
	v.OrderedPaths = append(v.OrderedPaths, path)
}

func (cg *contextGraph) processCurrentPath(echoCtx echo.Context) {
	path := echoCtx.Request().URL.Path

	for i, orderedPath := range cg.OrderedPaths {
		if orderedPath == path {
			cg.VertexPath = orderedPath

			if i < len(cg.OrderedPaths)-1 {
				cg.EdgePath = cg.OrderedPaths[i+1]
			} else {
				cg.EdgePath = ""
			}
			fmt.Printf("Vertex: %s, Edge: %s\n", cg.VertexPath, cg.EdgePath)
			return
		}
	}
}
