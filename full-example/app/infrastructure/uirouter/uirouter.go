package uirouter

import (
	"strings"

	"github.com/google/uuid"
)

type UIRouter struct {
	RootURL     string
	RootHTML    string
	activeRoute Route
	Routes      []Route
}

type Route struct {
	URL        string
	RedirectTo string
}

func (router UIRouter) GetActiveRoute(requestURL string) (Route, bool) {
	for _, route := range router.Routes {
		if route.URL == requestURL {
			if route.RedirectTo != "" {
				// Encuentra y devuelve la ruta a la que se redirige
				for _, redirectRoute := range router.Routes {
					if redirectRoute.URL == route.RedirectTo {
						return redirectRoute, true
					}
				}
			}
			return route, true
		}
	}
	return Route{}, false
}

func NewSelectorTarget() string {
	return "selector" + strings.ReplaceAll(uuid.NewString(), "-", "")
}
