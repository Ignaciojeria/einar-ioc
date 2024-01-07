# Golang Dependency Injection Framework 🪡

## 🔧 Installation
To install ioc, use the following command:

    go get github.com/Ignaciojeria/einar-ioc@v1.6.0


## 🔍 Tutorial : Before start

The dependencies in the framework are represented as a Directed Acyclic Graph (DAG) where:

1. Dependencies are regitered using ioc.Registry(vertex,...edges) function.
2. Constructors are vertices.
3. Dependencies are edges.
4. The graph is topologically ordered.
5. Dependency loading starts from descendant nodes and proceeds to their ancestors.

## 👨‍💻 Setup

As a first step, we'll make sure that the `main` function loads all the dependencies we will inject later on. This initial loading of dependencies is crucial for setting up our Dependency Injection framework. 

```go
package main

import (
	"os"
	ioc "github.com/Ignaciojeria/einar-ioc"
)
func main() {
	if err := ioc.LoadDependencies(); err != nil {
		os.Exit(1)
	}
}
```

## 👨‍💻 HTTP Router Registration
Here, we register the HTTP router using the einar-ioc framework. The NewRouter function is registered as a vertex in the dependency graph. This means NewRouter will be used to instantiate the Echo HTTP router when needed in the application.

```go
package router

import (
	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewRouter)

func NewRouter() *echo.Echo {
	e := echo.New()
	return e
}
```

#### 🔍 Retrieving the registered router & start the server
```go
package main

import (
	"log"
	"tutorial/app/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	// Retrieving the registered router & start the server
	r, _ := ioc.Get(router.NewRouter)
	r.(*echo.Echo).Start(":8080")
}
```

## 👨‍💻 HTTP Handler Registration
```go
package handler

import (
	"net/http"
	"tutorial/app/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

// Registers newGetExampleHandler as a constructor in the dependency injection container.
// It depends on router.NewRouter
var _ = ioc.Registry(newGetExampleHandler, router.NewRouter)

type getExampleHandler struct {
}

// newGetExampleHandler is a constructor function for getExampleHandler.
// It takes a *echo.Echo instance (r) as a parameter (edge in the dependency graph),
// indicating that it relies on the Echo router (vertex) for its operation.
func newGetExampleHandler(r *echo.Echo) getExampleHandler {
	handler := getExampleHandler{}
	r.GET("/example", handler.handle)
	return handler
}

func (h getExampleHandler) handle(c echo.Context) error {
	return c.String(http.StatusOK, "Mira mamá, sin manos!")
}
```


## 📑 Ioc.Registry : Constructor Registration Rules

The Ioc.Registry(vertex, edges...) function in our Inversion of Control (IoC) system plays a critical role in managing and registering dependencies. This function is specifically designed to register constructors that meet certain criteria regarding the types of values they return.

#### 🔍 Return Type Constraints
It's important to note that Ioc.Registry can only register constructors that return up to two specific types:

1. Associated Structure of the Constructor: This is the primary type of the object or structure that the constructor is designed to create. It's a mandatory requirement for each constructor to return this type.

2. Error (Optional): Optionally, the constructor can return a second type, which is an error. This return is used to indicate if there was any error during the creation of the object or structure. The inclusion of this return type is optional, but when present, it provides a robust way to handle errors in the creation process.