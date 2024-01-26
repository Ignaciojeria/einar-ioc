package uirouter

type UIRouter struct {
	RootURL     string
	RootHTML    string
	activeRoute Route
	Routes      []Route
}

type Route struct {
	URL        string
	RedirectTo string
	Children   []Child
}

type Child struct {
	Route
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

		// Revisar las rutas hijas
		for _, child := range route.Children {
			if child.URL == requestURL {
				return child.Route, true
			}
		}
	}
	return Route{}, false
}
