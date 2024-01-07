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

####  🔍 Create your router file
Navigate to the /app/router folder. Inside this folder, we will create router.go
```bash
/app/router
 - router.go #Echo Router 
``` 

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
In this section, we register the HTTP Handler using the einar-ioc framework. The `newGetExampleHandler` function is registered as a vertex in the dependency graph. Additionally, `router.NewRouter` is specified as a dependency of `newGetExampleHandler` and is represented as an edge in the graph. This setup signifies that `router.NewRouter` will be injected as a dependency into `newGetExampleHandler` during the instantiation of the handler.

####  🔍 Create your first handler
Navigate to the /app/handler folder. Inside this folder, we will create get_example_handler.go
```bash
/app/handler
 - get_example_handler.go
``` 

```go
package handler

import (
	"net/http"
	"tutorial/app/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(newGetExampleHandler, router.NewRouter)

type getExampleHandler struct {
}

func newGetExampleHandler(r *echo.Echo) getExampleHandler {
	handler := getExampleHandler{}
	r.GET("/example", handler.handle)
	return handler
}

func (h getExampleHandler) handle(c echo.Context) error {
	return c.String(http.StatusOK, "Mira mamá, sin manos!")
}
```

#### 🔍 Handler discovering
```go
package main

import (
	"log"
	_ "tutorial/app/handler" //Required to discover the handler and all its descendants.
	"tutorial/app/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	r, _ := ioc.Get(router.NewRouter)
	r.(*echo.Echo).Start(":8080")
}
```
#### 🔍 Start application and check Handler

<div align="center">
    <img src="/sketching/hello_world.jpeg" width="300" height="200">
</div>

## 📑 Ioc.Registry : Constructor Registration Rules

The Ioc.Registry(vertex, edges...) function in our Inversion of Control (IoC) system plays a critical role in managing and registering dependencies. This function is specifically designed to register constructors that meet certain criteria regarding the types of values they return.

#### 🔍 Return Type Constraints
It's important to note that Ioc.Registry can only register constructors that return up to two specific types:

1. Associated Structure of the Constructor: This is the primary type of the object or structure that the constructor is designed to create. It's a mandatory requirement for each constructor to return this type.

2. Error (Optional): Optionally, the constructor can return a second type, which is an error. This return is used to indicate if there was any error during the creation of the object or structure. The inclusion of this return type is optional, but when present, it provides a robust way to handle errors in the creation process.