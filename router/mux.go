package router

import (
	"refactor/ioc"

	"github.com/go-chi/chi/v5"
)

var _ = ioc.Installation(NewRouter)

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
