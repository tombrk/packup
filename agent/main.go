package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/structs"
	cron3 "github.com/robfig/cron/v3"

	"github.com/sh0rez/packup/internal/config"
	"github.com/sh0rez/packup/internal/logs"
	"github.com/sh0rez/packup/internal/restic"
)

var log = logs.Logger()
var Version = "dev"

func main() {
	exec := flag.Bool("x", false, "exec src and read dir from stdout")
	version := flag.Bool("version", false, "print version")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-x] <src> <cron>\n\nFlags:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	logs.Verbose(*verbose)

	if *version {
		fmt.Println(Version)
		return
	}

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := parse()
	if err != nil {
		log.Fatal().Err(err).Msgf("parsing config")
	}

	cfg.Src = config.Source{Path: flag.Arg(0)}
	if *exec {
		cfg.Src = config.Source{Exec: flag.Arg(0)}
	}

	ctx := context.Background()
	if err := run(ctx, cfg); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func run(ctx context.Context, cfg Config) error {
	repo, err := restic.Open(cfg.Repo, cfg.Password)
	if err != nil {
		return err
	}

	task := task{
		repo: *repo,
		src:  cfg.Src,
		sig:  make(chan string),
	}

	if cfg.Cron == "@now" {
		if ok := task.run("once"); !ok {
			os.Exit(1)
		}
		return nil
	}

	if err := cron(cfg.Cron, task.sig); err != nil {
		return err
	}

	stop := notify(task.sig, syscall.SIGHUP)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case trigger := <-task.sig:
			task.run(trigger)
		}
	}
}

type Config struct {
	Repo     string `env:"RESTIC_REPOSITORY"`
	Password string `env:"RESTIC_PASSWORD"`

	Src  config.Source
	Cron string `arg:"1"`
}

func parse() (Config, error) {
	var cfg Config
	var errs error

	for _, f := range structs.Fields(&cfg) {
		if key := f.Tag("env"); key != "" {
			val := os.Getenv(key)
			if val == "" {
				errs = errors.Join(errs, fmt.Errorf("env-var %s must be set", key))
			}
			f.Set(val)
		}

		if pos := f.Tag("arg"); pos != "" {
			i, err := strconv.Atoi(pos)
			if err != nil {
				panic(err)
			}
			val := flag.Arg(i)
			if val == "" {
				errs = errors.Join(errs, fmt.Errorf("pos-arg %d must be set", i))
			}
			f.Set(flag.Arg(i))
		}
	}

	return cfg, errs
}

type task struct {
	repo restic.Repository
	src  config.Source
	sig  chan string
}

func (t task) run(trigger string) bool {
	start := time.Now()
	dir, err := t.src.Dir()
	if err != nil {
		log.Err(err).Msg("failed reading backup source")
		return false
	}

	log.Info().Str("path", dir).Str("trigger", trigger).Msg("starting backup")
	result, err := t.repo.Backup(dir)
	if err != nil {
		log.Err(err).Dur("took", time.Since(start)).Msg("backup failed")
		return false
	}

	event := log.Info().Dur("took", time.Since(start))
	if result != nil && result.Summary != nil {
		event = event.
			Str("snapshot", result.Summary.SnapshotID).
			Uint("files", result.Summary.TotalFilesProcessed).
			Uint64("bytes", result.Summary.TotalBytesProcessed).
			Uint64("data_added", result.Summary.DataAdded).
			Uint64("data_added_packed", result.Summary.DataAddedPacked)
		if result.Summary.SnapshotID == "" {
			event = event.Bool("snapshot_skipped", true)
		}
	}
	event.Msg("backup finished")
	return true
}

func notify(sig chan<- string, signals ...os.Signal) (stop func()) {
	on := make(chan os.Signal)
	go func() {
		for s := range on {
			trigger(sig, fmt.Sprintf("signal(%s)", s))
		}
	}()
	signal.Notify(on, signals...)
	return func() { close(on) }
}

func cron(expr string, sig chan<- string) error {
	c := cron3.New()
	_, err := c.AddFunc(expr, func() {
		trigger(sig, fmt.Sprintf("cron(%s)", expr))
	})
	if err == nil {
		log.Info().Str("cron", expr).Msg("running on schedule")
		go c.Run()
	}
	return err
}

func trigger(sig chan<- string, from string) {
	select {
	case sig <- from:
	default:
		log.Error().Str("trigger", from).Msg("backup already in progress")
	}
}
