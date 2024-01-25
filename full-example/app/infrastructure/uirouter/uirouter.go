package uirouter

type UIRouter struct {
	Root   string
	Routes []Route
}

type Route struct {
	URL        string
	RedirectTo string
	Children   []Child
}

type Child struct {
	Route
}
