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
	home home.Home) map[string]uirouter.UIRouter {
	routerMap := map[string]uirouter.UIRouter{
		index.HTML: {
			RootHTML: index.HTML,
			Routes: []uirouter.Route{
				{
					//index.html router-outlet
					URL:        index.URL,
					RedirectTo: app.URL,
				},
				{
					//index.html router-outlet
					URL:        app.URL + home.URL,
					RedirectTo: app.URL,
				},
				{
					//index.html router-outlet required for RedirectTo
					URL: app.URL,
				},
				//index.html router-outlet
				{
					URL: home.URL,
				},
			},
		},
		app.HTML: {
			RootHTML: app.HTML,
			Routes: []uirouter.Route{
				{
					//app.html router-outlet
					URL:        app.URL,
					RedirectTo: app.URL + home.URL,
				},
				{
					//app.html router-outlet required for redirectTo
					URL: app.URL + home.URL,
				},
			},
		},
	}
	htmx.SetUIRouterMap(routerMap)
	server.Router().Use(htmx.middleware)
	server.Router().GET(app.URL, app.Render)
	server.Router().GET(app.URL+home.URL, home.Render)
	server.Router().GET(home.URL, home.Render)
	return routerMap
}
