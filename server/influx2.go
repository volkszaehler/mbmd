package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
)

// Influx2 is an InfluxDB v2 publisher
type Influx2 struct {
	sync.Mutex
	client      *influxdb.Client
	points      []*influxdb.RowMetric
	interval    time.Duration
	bucket      string
	org         string
	measurement string
}

// NewInflux2Client creates new publisher for influx
func NewInflux2Client(
	url string,
	bucket string,
	org string,
	measurement string,
	interval time.Duration,
	token string,
	user string,
	password string,
) *Influx2 {
	http := &http.Client{Timeout: writeTimeout}
	options := []influxdb.Option{
		influxdb.WithHTTPClient(http),
	}
	if token == "" {
		options = append(options, influxdb.WithUserAndPass(user, password))
	}

	client, err := influxdb.New(url, token, options...)
	if err != nil {
		log.Fatalf("influx2: error creating client: %v", err)
	}

	if bucket == "" {
		log.Fatal("influx2: missing bucket")
	}
	if measurement == "" {
		log.Fatal("influx2: missing measurement")
	}

	return &Influx2{
		client:      client,
		interval:    interval,
		measurement: measurement,
		bucket:      bucket,
		org:         org,
	}
}

// writeBatchPoints asynchronously writes the collected points to influx
func (m *Influx2) writeBatchPoints() {
	m.Lock()

	// get current batch
	if len(m.points) == 0 {
		m.Unlock()
		return
	}

	// replace current batch
	points := m.points
	m.points = nil
	m.Unlock()

	// write batch
	metrics := make([]influxdb.Metric, len(points))
	for i, p := range points {
		metrics[i] = p
	}

	if _, err := m.client.Write(context.Background(), m.bucket, m.org, metrics...); err != nil {
		log.Printf("influx2: failed writing %d points, will retry: %v", len(points), err)

		// put points back at beginning of next batch
		m.Lock()
		m.points = append(points, m.points...)
		m.Unlock()
	}
}

// asyncWriter periodically calls writeBatchPoints
func (m *Influx2) asyncWriter(exit <-chan bool) <-chan bool {
	done := make(chan bool) // signal writer stopped

	// async batch writer
	go func() {
		ticker := time.NewTicker(m.interval)
		for {
			select {
			case <-ticker.C:
				m.writeBatchPoints()
			case <-exit:
				ticker.Stop()
				m.writeBatchPoints()
				done <- true
				return
			}
		}
	}()

	return done
}

// Run Influx publisher
func (m *Influx2) Run(in <-chan QuerySnip) {
	// run async writer
	exit := make(chan bool)     // exit signals to stop writer
	done := m.asyncWriter(exit) // done signals writer stopped

	for snip := range in {
		p := influxdb.NewRowMetric(
			map[string]interface{}{"value": snip.Value},
			m.measurement,
			map[string]string{
				"device": snip.Device,
				"type":   snip.Measurement.String(),
			},
			snip.Timestamp,
		)

		m.Lock()
		m.points = append(m.points, p)
		m.Unlock()
	}

	// close write loop
	exit <- true
	<-done

	m.client.Close()
}
