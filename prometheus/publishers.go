package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// publisherCollectors contains all Prometheus metrics about publishers like HTTP, MQTT, ...
//
// Implements collectable interface
type publisherCollectors struct{}

var publisherMetricsLabels = []string{"type"}

var (
	PublisherCreated = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_created_total",
			"Total count of publishers created for publishing smart meter measurement data.",
		),
		publisherMetricsLabels,
	)

	PublisherDataPublished = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_data_published_total",
			"Total count of publish processes where smart meter measurement data is published.",
		),
		publisherMetricsLabels,
	)

	PublisherDataPublishAttempt = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_data_publish_attempts_total",
			"Total count of publish attempts where smart meter measurement data is published.",
		),
		publisherMetricsLabels,
	)

	PublisherDataPublishedSize = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_data_published_size_bytes_total",
			"Total bytes of sent smart meter measurement data via publishers.",
		),
		publisherMetricsLabels,
	)

	PublisherDataPublishedError = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_data_published_errors_total",
			"Total failures of publish processes where smart meter measurement is published.",
		),
		publisherMetricsLabels,
	)

	PublisherConnectionSuccess = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_connection_successes_total",
			"Total successful connections to external databases/data storages via publishers.",
		),
		publisherMetricsLabels,
	)

	PublisherConnectionFailure = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_connection_failures_total",
			"Total failed connections to external databases/data storages via publishers.",
		),
		publisherMetricsLabels,
	)

	PublisherConnectionFlush = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_connection_flushes_total",
			"Total connection flushes to external databases/data storages via publishers. Flushing equals write operation to external database/storage.",
		),
		publisherMetricsLabels,
	)

	PublisherConnectionTimeOut = prometheus.NewCounterVec(
		newCounterOpts(
			"publisher_connection_timeouts_total",
			"Total connection timeouts during publish/flush processes.",
		),
		publisherMetricsLabels,
	)
)

func (publisherCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		PublisherCreated,
		PublisherDataPublished,
		PublisherDataPublishedSize,
		PublisherDataPublishedError,
		PublisherDataPublishAttempt,
		PublisherConnectionSuccess,
		PublisherConnectionFlush,
		PublisherConnectionFailure,
		PublisherConnectionTimeOut,
	}
}
