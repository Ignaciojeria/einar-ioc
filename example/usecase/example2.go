package usecase

import (
	"refactor/ioc"
)

var _ = ioc.Registry(NewExample2)

type Example2 struct {
}

func NewExample2() Example2 {
	return Example2{}
}

func (e Example2) Execute() string {
	return "hello mom"
}
