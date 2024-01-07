package handler

import (
	"net/http"
	"tutorial/app/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

// Registers newGetExampleHandler as a constructor in the dependency injection container.
// It depends on router.NewRouter
var _ = ioc.Registry(newGetExampleHandler, router.NewRouter)

type getExampleHandler struct {
}

// newGetExampleHandler is a constructor function for getExampleHandler.
// It takes a *echo.Echo instance (r) as a parameter (edge in the dependency graph),
// indicating that it relies on the Echo router (vertex) for its operation.
func newGetExampleHandler(r *echo.Echo) getExampleHandler {
	handler := getExampleHandler{}
	r.GET("/example", handler.handle)
	return handler
}

func (h getExampleHandler) handle(c echo.Context) error {
	return c.String(http.StatusOK, "Mira mamá, sin manos!")
}
