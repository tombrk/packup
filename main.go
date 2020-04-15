package main

import (
	"log"
	"net/http"
	"os"

	"github.com/markbates/pkger"
	"github.com/sh0rez/packup/pkg/api"
)

func main() {
	if os.Getenv("RESTIC_PASSWORD") == "" {
		log.Fatalln("RESTIC_PASSWORD unset. No repository connection possible")
	}
	if os.Getenv("RESTIC_REPOSITORY") == "" {
		log.Fatalln("RESTIC_REPOSITORY unset. No repository connection possible")
	}

	// register api
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.New()))

	// register ui
	http.Handle("/", http.FileServer(pkger.Dir("/ui/build")))

	log.Println("Listening on :2112")
	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Fatalln(err)
	}
}
