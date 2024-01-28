package uirouter

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UIRouter struct {
	RootURL     string
	RootHTML    string
	activeRoute Route
	Routes      []Route
}

type Route struct {
	UserInputPath string
	Render        func(e echo.Context) error
	URL           string
	RedirectTo    string
}

func (router UIRouter) GetActiveRoute(requestURL string) (Route, bool) {
	for _, route := range router.Routes {
		if isMatch(route.URL, requestURL) {
			if route.RedirectTo != "" {
				// Encuentra y devuelve la ruta a la que se redirige
				for _, redirectRoute := range router.Routes {
					if redirectRoute.URL == route.RedirectTo {
						return redirectRoute, true
					}
				}
			}

			// Construir la nueva URL con parámetros reemplazados
			replacedURL := replaceParams(route.URL, requestURL)
			route.URL = replacedURL
			return route, true
		}
	}
	return Route{}, false
}

// replaceParams toma la URL de la ruta y la URL de la solicitud,
// y reemplaza los parámetros en la URL de la ruta con los valores de la solicitud.
func replaceParams(routeURL, requestURL string) string {
	routeParts := strings.Split(routeURL, "/")
	requestParts := strings.Split(requestURL, "/")

	for i, part := range routeParts {
		if strings.HasPrefix(part, ":") && i < len(requestParts) {
			routeParts[i] = requestParts[i]
		}
	}

	return strings.Join(routeParts, "/")
}

// isMatch compara la URL de la ruta y la URL de la solicitud.
// Considera que hay una coincidencia incluso si las partes variables son diferentes.
func isMatch(routeURL, requestURL string) bool {
	routeParts := strings.Split(routeURL, "/")
	requestParts := strings.Split(requestURL, "/")

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i := range routeParts {
		if routeParts[i] != requestParts[i] && !strings.HasPrefix(routeParts[i], ":") {
			return false
		}
	}

	return true
}

func NewSelectorTarget() string {
	return "selector" + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (r Route) WithUserInputPath(userInputPath string) Route {
	// Dividir las URLs en partes para compararlas
	r.UserInputPath = userInputPath
	routeParts := strings.Split(r.URL, "/")
	userInputParts := strings.Split(userInputPath, "/")

	// Verificar que las dos URLs tengan el mismo número de partes
	if len(routeParts) != len(userInputParts) {
		return r
	}

	// Recorrer las partes y reemplazar los parámetros
	for i, part := range routeParts {
		if strings.HasPrefix(part, ":") && i < len(userInputParts) {
			routeParts[i] = userInputParts[i]
		}
	}

	// Reconstruir la URL con los parámetros reemplazados
	r.URL = strings.Join(routeParts, "/")
	return r
}
