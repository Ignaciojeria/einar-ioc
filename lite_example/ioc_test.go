package ioc

import (
	"fmt"
	"testing"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

func init() {
	ioc.Registry(NewMessage)
	ioc.Registry(NewGreeter, NewMessage)
	ioc.Registry(NewEvent, NewGreeter)
}

type Message string

func NewMessage() Message {
	return Message("Hi there!")
}

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) {
	fmt.Println(g.Greet())
}

func TestLoadDependencies(t *testing.T) {
	if err := ioc.LoadDependencies(); err != nil {
		t.Fail()
	}
}
