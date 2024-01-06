package usecase

import "refactor/ioc"

var _ = ioc.Registry(
	NewOtherUseCase,
	NewExampleUseCase,
	NewExampleUseCase,
	NewOtherUseCase2,
	NewExample2)

type OtherUseCase struct {
	e ExampleUseCase
	o OtherUseCase2
}

func NewOtherUseCase(e ExampleUseCase, e2 ExampleUseCase, o OtherUseCase2, e3 Example2) OtherUseCase {
	return OtherUseCase{}
}

func (e OtherUseCase) Execute() string {
	return e.Execute()
}

var _ = ioc.Registry(NewOtherUseCase2)

type OtherUseCase2 struct {
	e ExampleUseCase
}

func NewOtherUseCase2() OtherUseCase2 {
	return OtherUseCase2{}
}

func (e OtherUseCase2) Execute() string {
	return e.Execute()
}
