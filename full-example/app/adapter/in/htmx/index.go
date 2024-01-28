package htmx

import (
	"embed"
	"my-project-name/app/infrastructure/server"
	"my-project-name/app/infrastructure/uicomponent"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(
	NewIndex,
	server.NewServer)

type Index struct {
	uicomponent.Component
}

//go:embed *.html
var html embed.FS

//go:embed *.css
var css embed.FS

func NewIndex(
	server server.Server) (Index, error) {
	view := Index{
		Component: uicomponent.Component{
			URL:  "/",
			HTML: "index.html",
			CSS:  "index.css",
		},
	}
	if err := server.TemplateRegistry(css, view.CSS); err != nil {
		return Index{}, err
	}
	if err := server.TemplateRegistry(html, view.HTML); err != nil {
		return Index{}, err
	}
	return view, nil
}
