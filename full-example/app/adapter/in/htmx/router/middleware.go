package router

import (
	"fmt"
	"my-project-name/app/configuration"
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
	conf configuration.Conf
	htmx *htmxm.HTMX
	//uiRouter    uirouter.UIRouter
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

		h := m.htmx.
			NewHandler(c.Response().Writer, c.Request())

		if h.IsHxRequest() && c.Request().Header.Get("HTMX-View") == "" {
			return next(c)
		}

		requestPath := c.Request().URL.Path
		ext := strings.ToLower(filepath.Ext(requestPath))

		// Omitir middleware para archivos .html, .css, y .js
		if ext == ".html" || ext == ".css" || ext == ".js" {
			return next(c)
		}
		fmt.Print("loop")

		view := c.Request().Header.Get("HTMX-View")

		if view == "" {
			view = "index.html"
		}

		router := m.uiRouterMap[view]
		activeRoute, found := router.GetActiveRoute(requestPath)
		if !found {
			return c.JSON(http.StatusNotFound, "Route not found")
		}
		if !h.IsHxRequest() {
			err := c.Render(http.StatusOK, "index.html", activeRoute)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
		c.Set("activeRoute", activeRoute)
		return next(c)
	}
}

func (s *middleware) SetUIRouterMap(u map[string]uirouter.UIRouter) {
	s.uiRouterMap = u
}
