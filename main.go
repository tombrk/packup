package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-clix/cli"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/sh0rez/packup/pkg/api"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/metrics"
	"github.com/sh0rez/packup/pkg/restic"
)

var Version string = "dev"

func main() {
	// setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

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
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
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

		// scheduler and job metrics
		c := cron.New()
		var scheduled []string
		watcher := make(metrics.RepoCollector)
		for name, job := range cfg.Jobs {
			job := job // <3 Go!
			log := log.With().Str("job", name).Logger()

			// repo metrics if locally available
			if fi, err := os.Stat(job.Repo); !os.IsNotExist(err) && fi.IsDir() {
				watcher[name] = job
			}

			// schedule backups if source defined
			if job.Source == "" {
				continue
			}

			rst := restic.New(job.Repo, job.Password, name)
			c.AddFunc(job.Schedule, func() {
				log.Info().Msg("Starting backup")
				if err := rst.Backup(job.Source, nil); err != nil {
					log.Error().Err(err).Msg("Backup failed")
				}
			})

			scheduled = append(scheduled, name)
		}

		c.Start()
		if len(scheduled) != 0 {
			log.Info().Strs("jobs", scheduled).Msg("Starting scheduler")
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
