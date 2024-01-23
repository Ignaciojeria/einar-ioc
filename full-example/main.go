package main

import (
	"log"
	"my-project-name/app/infrastructure/server"

	_ "my-project-name/app/adapter/in/controller"
	_ "my-project-name/app/adapter/in/htmx"
	_ "my-project-name/app/adapter/in/htmx/router"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	s, _ := ioc.Get(server.NewServer)
	s.(server.Server).Start()
}
