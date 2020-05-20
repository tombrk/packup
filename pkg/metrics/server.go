package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/restic"
	"github.com/sh0rez/packup/pkg/size"
)

const Namespace = "packup"

func WatchDir(name string, job config.Job) {
	labels := prometheus.Labels{"name": name}

	// directory sizes
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "repository_size_bytes",
			Help:        "Repository size in bytes on local disk",
			ConstLabels: labels,
		},
		func() float64 {
			n, err := size.Size(job.Repo)
			if err != nil {
				return -1
			}
			return float64(n)
		},
	)

	// snapshots count
	rst := restic.New(job.Repo, job.Password, name)
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "snapshots_total",
			Help:        "Count of snapshots in the repository",
			ConstLabels: labels,
		},
		func() float64 {
			s, err := rst.Snapshots()
			if err != nil {
				return -1
			}
			return float64(len(s))
		},
	)

	// last successful timestamp
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "last_snapshot_seconds",
			Help:        "UNIX seconds of the last successful snapshot",
			ConstLabels: labels,
		},
		func() float64 {
			s, err := rst.Snapshots()
			if err != nil {
				return -1
			}
			return float64(s[len(s)-1].Time.Unix())
		},
	)
}
