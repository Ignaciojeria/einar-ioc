package router

import (
	"my-project-name/app/adapter/in/htmx"
	"my-project-name/app/adapter/in/htmx/app"
	"my-project-name/app/adapter/in/htmx/home"
	"my-project-name/app/infrastructure/server"
	"my-project-name/app/infrastructure/uirouter"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(
	newUIRouter,
	NewMiddleware,
	server.NewServer,
	htmx.NewIndex,
	app.NewApp,
	home.NewHome)

type appRouter struct {
	index htmx.Index
	app   app.App
	home  home.Home
}

func newUIRouter(
	htmx middleware,
	server server.Server,
	index htmx.Index,
	app app.App,
	home home.Home) uirouter.UIRouter {
	router := uirouter.UIRouter{
		RootHTML: index.HTML,
		Routes: []uirouter.Route{
			{
				//index.html router-outlet
				URL:        index.URL,
				RedirectTo: app.URL,
			},
			{
				//index.html router-outlet
				URL: app.URL,
			},
			{
				//index.html router-outlet
				URL: home.URL,
			},
		},
	}
	htmx.SetUIRouter(router)
	server.Router().Use(htmx.middleware)
	server.Router().GET(app.URL, app.Render)
	server.Router().GET(home.URL, home.Render)
	return router
}
