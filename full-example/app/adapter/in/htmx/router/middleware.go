package router

import (
	"my-project-name/app/configuration"
	"my-project-name/app/infrastructure/uicomponent"
	"my-project-name/app/infrastructure/uirouter"
	"net/http"
	"path/filepath"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc"
	htmxm "github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewMiddleware, configuration.NewConf)

type middleware struct {
	conf        configuration.Conf
	htmx        *htmxm.HTMX
	uiRouterMap map[string]uirouter.UIRouter
}

func NewMiddleware(conf configuration.Conf) middleware {
	return middleware{
		htmx: htmxm.New(),
		conf: conf,
	}
}

func (m middleware) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		requestPath := c.Request().URL.Path
		if strings.HasPrefix(requestPath, m.conf.ApiPrefix) {
			return next(c)
		}

		ext := strings.ToLower(filepath.Ext(requestPath))
		if ext == ".html" || ext == ".css" || ext == ".js" {
			return next(c)
		}

		h := m.htmx.
			NewHandler(c.Response().Writer, c.Request())

		if h.IsHxRequest() && c.Request().Header.Get("HTMX-View") == "" {
			return next(c)
		}

		view := c.Request().Header.Get("HTMX-View")

		if view == "" {
			view = "index.html"
		}

		router := m.uiRouterMap[view]
		activeRoute, found := router.GetActiveRoute(requestPath)
		activeRoute = activeRoute.WithUserInputPath(c.Request().Header.Get("UserInputPath"))
		if !found {
			return c.JSON(http.StatusNotFound, "Route not found")
		}
		if !h.IsHxRequest() {
			activeRoute.UserInputPath = requestPath
			err := c.Render(http.StatusOK, "index.html", activeRoute)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
		c.Set(uicomponent.ActiveRoute, activeRoute)
		return next(c)
	}
}

func (s *middleware) SetUIRouterMap(u map[string]uirouter.UIRouter) {
	s.uiRouterMap = u
}
