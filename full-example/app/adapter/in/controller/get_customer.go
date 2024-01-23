package controller

import (
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newGetCustomer,
	server.NewServer)

type getCustomer struct {
	s server.Server
}

func newGetCustomer(s server.Server) getCustomer {
	controller := getCustomer{
		s: s,
	}
	controller.s.Router().GET(s.ApiPrefix()+"customer", controller.handle)
	return controller
}

func (ctrl getCustomer) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
