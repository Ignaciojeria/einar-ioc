package main

import (
	"fmt"
	"log"
	"net/http"
	_ "refactor/controller"
	"refactor/ioc"
	"refactor/router"
	_ "refactor/router"
	_ "refactor/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	router, err := ioc.Get(router.NewRouter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting server at port 3000")
	http.ListenAndServe(":3000", router.(*chi.Mux))
}
