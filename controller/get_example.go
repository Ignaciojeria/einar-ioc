package controller

import (
	"net/http"
	"refactor/ioc"
	"refactor/router"

	"github.com/go-chi/chi/v5"
)

// That should be inbound -__-
var _ = ioc.OutBoundAdapter(newGetExample, router.NewRouter)

type getExample struct {
}

func newGetExample(m *chi.Mux) getExample {
	c := getExample{}
	m.Get("/hello", c.handle)
	return c
}

func (c getExample) handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello mom"))
}
