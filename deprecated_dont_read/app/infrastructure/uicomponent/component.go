package uicomponent

import (
	"my-project-name/app/infrastructure/uirouter"

	"github.com/labstack/echo/v4"
)

const ActiveRoute = "activeRoute"

type Component struct {
	ActiveRoute uirouter.Route
	Target      string
	URL         string
	HTML        string
	CSS         string
}

func (c Component) WithContext(ctx echo.Context) Component {
	c.ActiveRoute = ctx.Get(ActiveRoute).(uirouter.Route)
	return c
}
