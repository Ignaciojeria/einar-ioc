package controller

import (
	"my-project-name/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPostCustomer,
	server.NewServer)

type postCustomer struct {
	s server.Server
}

func newPostCustomer(s server.Server) postCustomer {
	controller := postCustomer{
		s: s,
	}
	controller.s.Router().POST(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl postCustomer) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
