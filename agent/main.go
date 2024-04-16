package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-clix/cli"
	cron3 "github.com/robfig/cron/v3"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/logs"
	"github.com/sh0rez/packup/pkg/restic"
)

var log = logs.Logger()
var Version = "dev"

func main() {
	cmd := &cli.Command{
		Use:     "packup-agent [FLAGS] <SOURCE> <CRON>",
		Short:   "backup agent using restic",
		Version: Version,
		Args:    cli.ArgsExact(2),
	}

	verbose := cmd.Flags().BoolP("verbose", "v", false, "print debug info")
	exec := cmd.Flags().BoolP("exec", "x", false, "use SOURCE as a program that prints a path to stdout")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		logs.Verbose(*verbose)

		repo := os.Getenv("RESTIC_REPOSITORY")
		if repo == "" {
			return fmt.Errorf("RESTIC_REPOSITORY must be set")
		}

		pass := os.Getenv("RESTIC_PASSWORD")
		if pass == "" {
			return fmt.Errorf("RESTIC_PASSWORD must be set")
		}

		expr := args[1]
		src := config.Source{Path: args[0]}
		if *exec {
			src = config.Source{Exec: args[0]}
		}

		rst, err := restic.Open(repo, pass)
		if err != nil {
			return err
		}

		task := task{
			rst: *rst,
			src: src,
			sig: make(chan string),
		}

		if expr == "@now" {
			if ok := task.run("once"); !ok {
				os.Exit(1)
			}
			return nil
		}

		if err := cron(expr, task.sig); err != nil {
			return err
		}

		stop := notify(task.sig, syscall.SIGHUP)
		defer stop()

		ctx := context.Background()

		for {
			select {
			case <-ctx.Done():
				return nil
			case tr := <-task.sig:
				task.run(tr)
			}
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type task struct {
	rst restic.Repository
	src config.Source
	sig chan string
}

func (t task) run(trigger string) bool {
	start := time.Now()
	dir, err := t.src.Dir()
	if err != nil {
		log.Err(err).Msg("failed reading backup source")
		return false
	}

	log.Info().Str("path", dir).Str("trigger", trigger).Msg("starting backup")
	if err := t.rst.Backup(dir); err != nil {
		log.Err(err).Dur("took", time.Since(start)).Msg("backup failed")
		return false
	}

	log.Info().Dur("took", time.Since(start)).Msg("backup finished")
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
