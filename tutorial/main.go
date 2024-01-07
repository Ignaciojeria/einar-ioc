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
	r, _ := ioc.Get(router.NewRouter)
	r.(*echo.Echo).Start(":8080")
}
