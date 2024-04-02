# Golang Dependency Injection Framework 🪡

## 🔧 Installation
To install ioc, use the following command:

    go get github.com/Ignaciojeria/einar-ioc@v1.9.0

## 👨‍💻 Example

```go
package ioc

import (
	"fmt"
	"testing"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(NewMessage)

type Message string

func NewMessage() Message {
	return Message("Hi there!")
}

var _ = ioc.Registry(NewGreeter, NewMessage)

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

var _ = ioc.Registry(NewEvent, NewGreeter)

type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) error {
	fmt.Println(g.Greet())
	return nil
}

func TestLoadDependencies(t *testing.T) {
	if err := ioc.LoadDependencies(); err != nil {
		t.Log(err)
		t.Fail()
	}
}
```
