package router

import (
	"my-project-name/app/infrastructure/uirouter"

	ioc "github.com/Ignaciojeria/einar-ioc"
	htmxm "github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewMiddleware)

type middleware struct {
	htmx     *htmxm.HTMX
	uiRouter uirouter.UIRouter
}

func NewMiddleware() middleware {
	return middleware{
		htmx: htmxm.New(),
	}
}

func (m middleware) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := m.htmx.
			NewHandler(c.Response().Writer, c.Request())

		if !h.IsHxRequest() {
			//TODO :  Enviar en la redirección el contexto del nodo desde dónde se redireccionó
			//TODO :  Ver la manera de representar multiples grafos para tener diferentes modulos de ruteo
			//contextGraph := m.server.ContextGraph(c)
			//	fmt.Println("Solicitud htmx!", contextGraph)
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

		return next(c)
	}
}

func (s *middleware) SetUIRouter(u uirouter.UIRouter) {
	s.uiRouter = u
}
