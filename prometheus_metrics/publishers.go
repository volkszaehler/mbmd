package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

// publisherCollectors contains all Prometheus metrics about publishers like HTTP, MQTT, ...
//
// Implements collectable interface
type publisherCollectors struct {}

var labels = []string{"type"}

var (
	PublisherCreated = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_created_total",
			"",
		),
		labels,
	)

	PublisherDataPublished = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_total",
			"",
		),
		labels,
	)

	PublisherDataPublishedSize = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_size_bytes_total",
			"",
		),
		labels,
	)

	PublisherDataPublishedError = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_data_published_errors_total",
			"",
		),
		labels,
	)

	PublisherConnectionFlush = prometheus.NewCounterVec(
		*newCounterOpts(
			"publisher_connection_flushes_total",
			"",
		),
		labels,
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

