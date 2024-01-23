package router

import (
	"my-project-name/app/adapter/in/htmx"
	"my-project-name/app/adapter/in/htmx/app"
	"my-project-name/app/infrastructure/server"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(
	newAppRouter,
	server.NewServer,
	htmx.NewIndex,
	app.NewApp)

type appRouter struct {
	server server.Server
	index  htmx.Index
	app    app.App
}

func newAppRouter(
	s server.Server,
	index htmx.Index,
	app app.App) appRouter {
	indexVertex, _ := s.AppDag().AddVertex(index.URL)
	appVertex, _ := s.AppDag().AddVertex(app.URL)
	s.AppDag().AddEdge(indexVertex, appVertex)
	s.Router().GET(index.URL, index.Render)
	s.Router().GET(app.URL, app.Render)
	return appRouter{}
}
