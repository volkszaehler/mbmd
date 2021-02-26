package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

type publisherCollectors struct {
	collectable
}

var labels = []string{"type"}

var (
	PublisherCreated = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_created_total",
			"",
		),
		[]string{"type"},
	)

	PublisherDataPublished = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_total",
			"",
		),
		[]string{"type"},
	)

	PublisherDataPublishedSize = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_size_bytes_total",
			"",
		),
		[]string{"type"},
	)

	PublisherDataPublishedError = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_errors_total",
			"",
		),
		[]string{"type"},
	)

	PublisherConnectionFlush = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_connection_flushes_total",
			"",
		),
		[]string{"type"},
	)
)

func (c publisherCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		PublisherCreated,
		PublisherDataPublished,
		PublisherDataPublishedSize,
		PublisherDataPublishedError,
		PublisherConnectionFlush,
	}
}
