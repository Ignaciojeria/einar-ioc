package ioc

import (
	"fmt"
	"testing"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var container = ioc.New()

func init() {
	container.RegistryAtEnd(AtEnd, NewEvent)
	container.Registry(NewMessage)
	container.Registry(NewGreeter, NewMessage)
	container.Registry(NewEvent, NewGreeter)
}

func AtEnd(gr Event) {
	fmt.Println("hello at end : " + gr.Greeter.Greet())
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

func NewEvent(g Greeter) Event {
	fmt.Println(g.Greet())
	return Event{
		Greeter: g,
	}
}

func TestLoadDependencies(t *testing.T) {
	if err := container.LoadDependencies(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
