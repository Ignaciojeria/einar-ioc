# Golang Minimalist Dependency Injection Framework 🪡

## 🔧 Installation
To install ioc, use the following command:

    go get github.com/Ignaciojeria/einar-ioc/v2

## 👨‍💻 Example

```go
package main

import (
	"fmt"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	//dont worry about dependency order registration.
	ioc.Registry(NewMessage)
	ioc.Registry(NewGreeter, NewMessage)
	ioc.Registry(NewEvent, NewGreeter)
	ioc.Registry(SendGreetFromEvent, NewEvent)
	/* this works too
	ioc.Registry(SendGreetFromEvent, NewEvent)
	ioc.Registry(NewGreeter, NewMessage)
	ioc.Registry(NewEvent, NewGreeter)
	ioc.Registry(NewMessage)
	*/
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

func SendGreetFromEvent(e Event) {
	fmt.Println(e.Greeter.Greet())
}

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		os.Exit(0)
	}
}
```
