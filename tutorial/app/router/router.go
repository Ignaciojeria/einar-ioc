package router

import (
	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

// First arguments is a vertext
// that means that NewRouter constructor was registered as vertex
var _ = ioc.Registry(NewRouter)

func NewRouter() *echo.Echo {
	e := echo.New()
	return e
}
