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

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) SendGreet() Message {
	return e.Greeter.Greet()
}

func TestLoadDependencies(t *testing.T) {
	if err := ioc.LoadDependencies(); err != nil {
		t.Log(err)
		t.Fail()
	}
	event := ioc.Get[Event](NewEvent)
	fmt.Println(event.SendGreet())
}
