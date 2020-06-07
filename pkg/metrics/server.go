package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/restic"
)

const Namespace = "packup"

type RepoCollector struct {
	Name string
	Job  config.Job
}

func (r *RepoCollector) Describe(d chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(r, d)
}

var (
	snapshotsTotalDesc = prometheus.NewDesc(
		Namespace+"_snapshots_total",
		"Count of snapshots in the repository",
		[]string{"name"}, nil,
	)
	lastSnapshotSecondsDesc = prometheus.NewDesc(
		Namespace+"last_snapshot_seconds",
		"UNIX seconds of the last successful snapshot",
		[]string{"name"}, nil,
	)
)

func (r *RepoCollector) Collect(c chan<- prometheus.Metric) {
	rst := restic.New(r.Job.Repo, r.Job.Password, r.Name)
	log := log.With().Str("job", r.Name).Logger()

	s, err := rst.Snapshots()
	if err != nil {
		log.Error().Err(err).Msg("Listing snapshots failed")
	}

	c <- prometheus.MustNewConstMetric(snapshotsTotalDesc, prometheus.GaugeValue, float64(len(s)), r.Name)
	c <- prometheus.MustNewConstMetric(lastSnapshotSecondsDesc, prometheus.GaugeValue, float64(s[len(s)-1].Time.Unix()), r.Name)
}
