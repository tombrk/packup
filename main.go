package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-clix/cli"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sh0rez/packup/pkg/api"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/logs"
	"github.com/sh0rez/packup/pkg/metrics"
)

var Version string = "dev"

var log = logs.Logger()

func main() {
	cmd := &cli.Command{
		Use:     "packup-server",
		Short:   "Easy and efficient backups using Restic",
		Version: Version,
	}

	configFile := cmd.Flags().String("config", "packup.yaml", "YAML config file")
	verbose := cmd.Flags().BoolP("verbose", "v", false, "Print debug info")
	listen := cmd.Flags().String("listen", ":9763", "Network address to listen on")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if *verbose {
			logs.Verbose(true)
			log.Debug().Msg("Enabling debug logs")
		}

		cfg, err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal().Err(err).Str("file", *configFile).Msg("Loading config")
		}

		// register api
		if cfg.API || cfg.UI {
			log.Info().Msg("Handling API at /api/v1")
			http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.New(cfg.Jobs)))
		}

		// job metrics
		watcher := make(metrics.RepoCollector)
		for name, job := range cfg.Jobs {
			job := job // <3 Go!

			// repo metrics if locally available
			if fi, err := os.Stat(job.Repo); !os.IsNotExist(err) && fi.IsDir() {
				watcher[name] = job
			}
		}

		// register prometheus metrics
		prometheus.MustRegister(watcher)
		http.Handle("/metrics", promhttp.Handler())
		if len(watcher) != 0 {
			log.Info().Strs("jobs", watcher.Jobs()).Msg("Exposing repository metrics at /metrics")
		}

		// register ui (must be last)
		if cfg.UI {
			spa := &spaFs{pkged: http.FS(uiFs)}
			http.Handle("/", http.FileServer(spa))
			log.Info().Msg("Serving web-ui at /")
		}

		log.Info().Msgf("Listening on %s", *listen)
		if err := http.ListenAndServe(*listen, nil); err != nil {
			return err
		}

		return nil
	}

	err := cmd.Execute()
	if help := (cli.ErrHelp{}); errors.As(err, &help) {
		fmt.Println(help.Error())
	} else if err != nil {
		log.Error().Err(err).Msg("")
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
