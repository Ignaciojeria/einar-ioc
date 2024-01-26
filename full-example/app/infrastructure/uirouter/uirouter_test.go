package uirouter

import "testing"

func TestGetActiveRoute(t *testing.T) {

	router := UIRouter{
		RootHTML: "index.html",
		Routes: []Route{
			{
				URL:        "/",
				RedirectTo: "/app",
			},
			{
				URL: "/app",
			},
			{
				URL: "/home",
			},
		},
	}

	activeRoute, _ := router.GetActiveRoute("/")

	if activeRoute.URL != "/app" {
		t.Fail()
	}

	activeRoute, _ = router.GetActiveRoute("/app")

	if activeRoute.URL != "/home" {
		t.Fail()
	}

}
