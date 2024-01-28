package uirouter

import "testing"

func TestGetActiveRoute(t *testing.T) {

	router := UIRouter{
		RootHTML: "index.html",
		Routes: []Route{
			{
				//index.html router-outlet
				URL:        "/",
				RedirectTo: "/app",
			},
			{
				//index.html router-outlet
				URL: "/app",
			},
			{
				//index.html router-outlet
				URL: "/home",
			},
			{
				//index.html router-outlet
				URL: "/users/:id",
			},
			{
				//index.html router-outlet
				URL: "/users",
			},
			{
				//index.html router-outlet
				URL: "/customers/:id/company/:id",
			},
		},
	}

	activeRoute, _ := router.GetActiveRoute("/")

	if activeRoute.URL != "/app" {
		t.Fail()
	}

	activeRoute, _ = router.GetActiveRoute("/app")

	if activeRoute.URL != "/app" {
		t.Fail()
	}

	activeRoute, _ = router.GetActiveRoute("/users/1")

	if activeRoute.URL != "/users/1" {
		t.Fail()
	}

	activeRoute, _ = router.GetActiveRoute("/users")

	if activeRoute.URL != "/users" {
		t.Fail()
	}

	activeRoute, _ = router.GetActiveRoute("/customers/1/company/2")

	// me refiero que el objeto activeRoute.URL en vez de retornar lo siguiente :
	if activeRoute.URL != "/customers/1/company/2" {
		t.Fail()
	}

}
