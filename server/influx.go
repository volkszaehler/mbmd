package server

import (
	"fmt"
	"github.com/volkszaehler/mbmd/prometheus_metrics"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
	api "github.com/influxdata/influxdb-client-go/api"
)

// Influx is an InfluxDB v2 publisher
type Influx struct {
	client      influxdb.Client
	writer      api.WriteAPI
	measurement string
}

// NewInfluxClient creates new publisher for influx
func NewInfluxClient(
	url string,
	database string,
	measurement string,
	org string,
	token string,
	user string,
	password string,
) *Influx {
	// InfluxDB v1 compatibility
	if token == "" && user != "" {
		token = fmt.Sprintf("%s:%s", user, password)
	}

	client := influxdb.NewClient(url, token)

	if database == "" {
		log.Fatal("influx: missing database")
	}
	if measurement == "" {
		log.Fatal("influx: missing measurement")
	}

	prometheus_metrics.PublisherCreated.WithLabelValues("influx").Inc()

	return &Influx{
		client:      client,
		measurement: measurement,
		writer:      client.WriteApi(org, database),
	}
}

// Run Influx publisher
func (m *Influx) Run(in <-chan QuerySnip) {
	// log errors
	go func() {
		for err := range m.writer.Errors() {
			log.Printf("influxdb error: %v", err)
			prometheus_metrics.PublisherDataPublishedError.WithLabelValues("influx").Inc()
		}
	}()

	for snip := range in {
		tags := map[string]string{
			"device": snip.Device,
			"type":   snip.Measurement.String(),
		}

		fields := map[string]interface{}{
			"value": snip.Value,
		}

		// write asynchronously
		p := influxdb.NewPoint(m.measurement, tags, fields, time.Now())
		m.writer.WritePoint(p)
		// prometheus_metrics.PublisherDataPublishedSize.WithLabelValues("influx").Add(float64(len(p)))
		prometheus_metrics.PublisherDataPublished.WithLabelValues("influx").Inc()
	}

	m.client.Close()
	prometheus_metrics.PublisherConnectionFlush.WithLabelValues("influx").Inc()
}
