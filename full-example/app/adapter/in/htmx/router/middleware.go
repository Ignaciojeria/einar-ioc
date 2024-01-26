package router

import (
	"my-project-name/app/configuration"
	"my-project-name/app/infrastructure/uirouter"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	htmxm "github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewMiddleware, configuration.NewConf)

type middleware struct {
	conf     configuration.Conf
	htmx     *htmxm.HTMX
	uiRouter uirouter.UIRouter
}

func NewMiddleware(conf configuration.Conf) middleware {
	return middleware{
		htmx: htmxm.New(),
		conf: conf,
	}
}

func (m middleware) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := m.htmx.
			NewHandler(c.Response().Writer, c.Request())
		activeRoute, found := m.uiRouter.GetActiveRoute(c.Request().URL.Path)
		if !found {
			return c.JSON(http.StatusNotFound, "Route not found")
		}
		if !h.IsHxRequest() {
			err := c.Render(http.StatusOK, m.uiRouter.RootHTML, activeRoute)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
		c.Set("activeRoute", activeRoute)
		return next(c)
	}
}

func (s *middleware) SetUIRouter(u uirouter.UIRouter) {
	s.uiRouter = u
}
