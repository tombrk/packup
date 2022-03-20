package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-clix/cli"
	"github.com/robfig/cron/v3"
	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/logs"
	"github.com/sh0rez/packup/pkg/restic"
)

var log = logs.Logger()
var Version = "dev"

func main() {
	cmd := &cli.Command{
		Use:     "packup-agent [FLAGS] <SOURCE> <CRON>",
		Short:   "Backup agent using restic",
		Version: Version,
		Args:    cli.ArgsExact(2),
	}

	verbose := cmd.Flags().BoolP("verbose", "v", false, "Print debug info")
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

		schedule := args[1]
		source := config.Source{Path: args[0]}
		if *exec {
			source = config.Source{Exec: args[0]}
		}

		rst, err := restic.New(repo, pass)
		if err != nil {
			return err
		}

		c := cron.New()
		c.AddFunc(schedule, func() {
			start := time.Now()
			dir, err := source.Dir()
			if err != nil {
				log.Err(err).Msg("Failed reading backup source")
				return
			}

			log.Info().Str("path", dir).Msg("Starting backup")
			if err := rst.Backup(dir); err != nil {
				log.Err(err).Dur("took", time.Since(start)).Msg("Backup failed")
				return
			}

			log.Info().Dur("took", time.Since(start)).Msg("Backup finished")
		})

		log.Info().Str("schedule", schedule).Msg("Running")

		c.Run()
		return nil
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
