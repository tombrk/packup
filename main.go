package main

import (
	"log"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/sh0rez/packup/pkg/api"
)

func main() {
	// register api
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.New()))

	// register ui
	http.Handle("/", http.FileServer(pkger.Dir("/ui/build")))

	log.Println("Listening on :2112")
	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Fatalln(err)
	}
}
