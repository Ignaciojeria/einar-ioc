package usecase

import (
	"refactor/ioc"
	"refactor/router"

	"github.com/go-chi/chi/v5"
)

var _ = ioc.Registry(NewExampleUseCase, router.NewRouter)

type ExampleUseCase struct {
}

func NewExampleUseCase(m *chi.Mux) ExampleUseCase {
	return ExampleUseCase{}
}

func (e ExampleUseCase) Execute() string {
	return "hello mom"
}
