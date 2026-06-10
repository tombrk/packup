package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/rs/zerolog/log"
	"github.com/sh0rez/packup/internal/config"
	"github.com/sh0rez/packup/internal/restic"
)

const Namespace = "packup"

type RepoCollector map[string]config.Job

func (r RepoCollector) Describe(d chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(r, d)
}

var (
	snapshotsTotalDesc = prometheus.NewDesc(
		Namespace+"_snapshots_total",
		"Count of snapshots in the repository",
		[]string{"name"}, nil,
	)
	lastSnapshotSecondsDesc = prometheus.NewDesc(
		Namespace+"_last_snapshot_seconds",
		"UNIX seconds of the last successful snapshot",
		[]string{"name"}, nil,
	)
)

func (r RepoCollector) Collect(c chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	for name, job := range r {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log := log.With().Str("job", name).Logger()

			rst, err := restic.Open(job.Repo, job.Password)
			if err != nil {
				log.Error().Err(err).Msg("Opening repo failed")
				return
			}

			s, err := rst.Snapshots()
			if err != nil {
				log.Error().Err(err).Msg("Listing snapshots failed")
				return
			}

			c <- prometheus.MustNewConstMetric(snapshotsTotalDesc, prometheus.GaugeValue, float64(len(s)), name)
			c <- prometheus.MustNewConstMetric(lastSnapshotSecondsDesc, prometheus.GaugeValue, float64(s[len(s)-1].Time.Unix()), name)
		}()
		wg.Wait()
	}
}

func (r RepoCollector) Jobs() []string {
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	return keys
}
