package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/restic"
	"github.com/sh0rez/packup/pkg/size"
)

const Namespace = "packup"

func Server(jobs config.Jobs) http.Handler {
	for name, job := range jobs {
		job := job // <3 Go

		// directory sizes
		promauto.NewGaugeFunc(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Name:        "repository_size_bytes",
				Help:        "Repository size in bytes on local disk",
				ConstLabels: prometheus.Labels{"name": name},
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
		rst := restic.New(job.Repo, job.Password)
		promauto.NewGaugeFunc(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Name:        "snapshots_total",
				Help:        "Count of snapshots in the repository",
				ConstLabels: prometheus.Labels{"name": name},
			},
			func() float64 {
				s, err := rst.Snapshots()
				if err != nil {
					return -1
				}
				return float64(len(s))
			},
		)
	}

	return promhttp.Handler()
}
