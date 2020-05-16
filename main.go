package main

import (
	"log"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/sh0rez/packup/pkg/api"
	"github.com/sh0rez/packup/pkg/config"
)

func main() {
	cfg, err := config.LoadFile("packup.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	apiHandler, err := api.New(cfg.Jobs)
	if err != nil {
		log.Fatalln(err)
	}

	// register api
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", apiHandler))

	// register ui
	http.Handle("/", http.FileServer(pkger.Dir("/ui/build")))

	log.Println("Listening on :2112")
	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Fatalln(err)
	}
}
