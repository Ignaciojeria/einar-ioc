package router

import (
	"my-project-name/app/adapter/in/htmx"
	"my-project-name/app/adapter/in/htmx/app"
	"my-project-name/app/infrastructure/server"
	"my-project-name/app/infrastructure/uirouter"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(
	newUIRouter,
	NewMiddleware,
	server.NewServer,
	htmx.NewIndex,
	app.NewApp)

type appRouter struct {
	index htmx.Index
	app   app.App
}

func newUIRouter(
	htmx middleware,
	server server.Server,
	index htmx.Index,
	app app.App) appRouter {
	server.Router().Use(htmx.middleware)
	server.Router().GET(index.URL, index.Render)
	server.Router().GET(app.URL, app.Render)
	router := uirouter.UIRouter{
		Root: index.URL,
		Routes: []uirouter.Route{
			{
				URL:        index.URL,
				RedirectTo: app.URL,
			},
		},
	}
	htmx.SetUIRouter(router)
	return appRouter{}
}
