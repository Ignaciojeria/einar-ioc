# Golang Dependency Injection Framework🔥

## 🔧 Installation
To install ioc, use the following command:

    go get github.com/Ignaciojeria/einar-ioc@v1.5.0


## 🔍 Tutorial : Before start

The dependencies in the framework are represented as a Directed Acyclic Graph (DAG) where:

1. Constructors are vertices.
2. Dependencies are edges.
3. The graph is topologically ordered.
4. Dependency loading starts from dependent nodes and proceeds to their ancestors.

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

