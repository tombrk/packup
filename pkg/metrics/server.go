package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"

	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/size"
)

const Namespace = "packup"

func WatchDir(name string, job config.Job) {
	log := log.With().Str("job", name).Logger()
	labels := prometheus.Labels{"name": name}

	tick := time.NewTicker(30 * time.Second)
	mut := &sync.Mutex{}

	var sizeBytes float64
	go func() {
		first := true
		for {
			if !first {
				_ = <-tick.C
			}
			first = false

			n, err := size.Size(job.Repo)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to calculate repo size")
				continue
			}

			mut.Lock()
			sizeBytes = float64(n)
			mut.Unlock()

			log.Debug().Float64("size", sizeBytes).Msg("Repository size calculated")
		}
	}()

	// directory sizes
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "repository_size_bytes",
			Help:        "Repository size in bytes on local disk",
			ConstLabels: labels,
		},
		func() float64 {
			mut.Lock()
			defer mut.Unlock()

			return float64(sizeBytes)
		},
	)

	// snapshots count
	// rst := restic.New(job.Repo, job.Password, name)
	// promauto.NewGaugeFunc(
	// 	prometheus.GaugeOpts{
	// 		Namespace:   Namespace,
	// 		Name:        "snapshots_total",
	// 		Help:        "Count of snapshots in the repository",
	// 		ConstLabels: labels,
	// 	},
	// 	func() float64 {
	// 		t := time.Now()
	// 		s, err := rst.Snapshots()
	// 		if err != nil {
	// 			return -1
	// 		}
	// 		log.Debug().Str("job", name).Dur("took", time.Since(t)).Msg("snapshots_total")
	// 		return float64(len(s))
	// 	},
	// )

	// last successful timestamp
	// promauto.NewGaugeFunc(
	// 	prometheus.GaugeOpts{
	// 		Namespace:   Namespace,
	// 		Name:        "last_snapshot_seconds",
	// 		Help:        "UNIX seconds of the last successful snapshot",
	// 		ConstLabels: labels,
	// 	},
	// 	func() float64 {
	// 		s, err := rst.Snapshots()
	// 		t := time.Now()
	// 		if err != nil {
	// 			return -1
	// 		}
	// 		log.Debug().Str("job", name).Dur("took", time.Since(t)).Msg("last_snapshot_seconds")
	// 		return float64(s[len(s)-1].Time.Unix())
	// 	},
	// )
}
