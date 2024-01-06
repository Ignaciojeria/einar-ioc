package controller

import (
	"net/http"
	"refactor/ioc"
	"refactor/router"
	"refactor/usecase"

	"github.com/go-chi/chi/v5"
)

var _ = ioc.Registry(
	newGetExample,
	router.NewRouter,
	usecase.NewExampleUseCase,
	usecase.NewOtherUseCase)

type getExample struct {
	u usecase.ExampleUseCase
	o usecase.OtherUseCase
}

func newGetExample(m *chi.Mux, u usecase.ExampleUseCase, o usecase.OtherUseCase) getExample {
	c := getExample{
		u: u,
		o: o,
	}
	m.Get("/hello", c.handle)
	return c
}

func (c getExample) handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.u.Execute()))
}
