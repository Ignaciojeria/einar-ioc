package usecase

//var _ = ioc.Registry(ioc.UseCase, NewExampleUseCase)

type ExampleUseCase struct {
}

func NewExampleUseCase() ExampleUseCase {
	return ExampleUseCase{}
}
