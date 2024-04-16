package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sh0rez/packup/pkg/api"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/logs"
	"github.com/sh0rez/packup/pkg/metrics"
	"github.com/sh0rez/packup/ui"
)

var Version string = "dev"
var log = logs.Logger()

func main() {
	file := flag.String("config", "packup.yaml", "yaml config file")
	listen := flag.String("listen", ":9763", "http address to listen on")
	verbose := flag.Bool("v", false, "debug logging")
	version := flag.Bool("version", false, "")
	flag.Parse()
	logs.Verbose(*verbose)

	if *version {
		fmt.Println(Version)
		return
	}

	cfg, err := config.Load(*file)
	if err != nil {
		log.Fatal().Err(err).Str("file", *file).Msg("loading config")
	}

	ctx := context.Background()
	if err := run(ctx, *cfg, *listen); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func run(ctx context.Context, cfg config.Config, addr string) error {
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.New(cfg.Jobs)))

	col := make(metrics.RepoCollector)
	for name, job := range cfg.Jobs {
		if fi, err := os.Stat(job.Repo); !os.IsNotExist(err) && fi.IsDir() {
			col[name] = job
		}
	}
	prometheus.MustRegister(col)
	http.Handle("/metrics", promhttp.Handler())

	ui := appfs{FileSystem: http.FS(ui.Files())}
	http.Handle("/", http.FileServer(ui))

	log.Info().Msgf("listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// appfs implements http.FileSystem to serve a single page application,
// returning index.html if the named file does not exist
type appfs struct {
	http.FileSystem
}

func (fs appfs) Open(name string) (http.File, error) {
	if f, err := fs.FileSystem.Open(name); err == nil {
		return f, nil
	}
	return fs.FileSystem.Open("index.html")
}
