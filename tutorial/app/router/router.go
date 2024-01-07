package router

import (
	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewRouter)

func NewRouter() *echo.Echo {
	e := echo.New()
	return e
}
