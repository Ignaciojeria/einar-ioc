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
	//print: Hi there!
	fmt.Println(ioc.Get[Message](NewMessage))
	mb := ioc.NewMockBehaviourForTesting[Message](NewMessage, Message("message mocked"))
	//print: message mocked
	fmt.Println(ioc.Get[Message](NewMessage))
	mb.Release()
	mb2 := ioc.NewMockBehaviourForTesting[Message](NewMessage, Message("message mocked 2"))
	//print: message mocked 2
	fmt.Println(ioc.Get[Message](NewMessage))
	mb2.Release()
	//print: Hi there!
	fmt.Println(ioc.Get[Message](NewMessage))
}
