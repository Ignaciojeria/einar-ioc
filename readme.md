# Golang Minimalist Dependency Injection Framework 🪡

## 🔧 Installation
To install ioc, use the following command:

    go get github.com/Ignaciojeria/einar-ioc/v2@v2.4.0

## 👨‍💻 Example

```go
package ioc

import (
	"fmt"
	"testing"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

func init() {
	// No need to worry about the order in which dependencies are registered here,
	// the framework will resolve them in the correct topological order.
	ioc.RegistryAtEnd(AtEnd, NewEvent)
	ioc.Registry(NewMessage)
	ioc.Registry(NewGreeter, NewMessage)
	ioc.Registry(NewEvent, NewGreeter)
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
	if err := ioc.LoadDependencies(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
```
