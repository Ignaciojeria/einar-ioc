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

}
