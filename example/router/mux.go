package router

import (
	"refactor/ioc"

	"github.com/go-chi/chi/v5"
)

var _ = ioc.Registry(NewRouter)

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
