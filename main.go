package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-clix/cli"
	"github.com/markbates/pkger"

	"github.com/sh0rez/packup/pkg/api"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/metrics"
)

var Version string = "dev"

func main() {
	cmd := &cli.Command{
		Use:     "packup",
		Short:   "Easy and efficient backups using Restic",
		Version: Version,
	}

	configFile := cmd.Flags().String("config", "packup.yaml", "YAML config file")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		cfg, err := config.LoadFile(*configFile)
		if err != nil {
			return err
		}

		apiHandler, err := api.New(cfg.Jobs)
		if err != nil {
			return err
		}

		// register api
		http.Handle("/api/v1/", http.StripPrefix("/api/v1", apiHandler))

		// register metrics
		http.Handle("/metrics", metrics.Server(cfg.Jobs))

		// register ui
		spa := &spaFs{pkged: pkger.Dir("/ui/build")}
		http.Handle("/", http.FileServer(spa))

		log.Println("Listening on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			return err
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// spaFs implements http.FileSystem to serve a single page application, meaning
// it always returns `index.html`
type spaFs struct {
	pkged http.FileSystem
}

func (fs *spaFs) Open(name string) (http.File, error) {
	f, err := fs.pkged.Open(name)
	if os.IsNotExist(err) {
		return fs.pkged.Open("/index.html")
	}
	return f, err
}
