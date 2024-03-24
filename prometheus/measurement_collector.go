package prometheus

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MeasurementCounterCollector is a struct which takes care of collecting prometheus.Counter metrics
// Whenever prometheus.Collect is called, all values of MeasurementCounterCollector.values are flushed
// and sent to the collector channel.
type MeasurementCounterCollector struct {
	metricsMap *prometheus.MetricVec
	desc       *prometheus.Desc
	values     struct {
		sync.RWMutex
		data map[string]*measurementResult // Contains latest value of meters.MeasurementResult
	}
	fqName string
	opts   prometheus.CounterOpts
}

func NewMeasurementCounterCollector(opts prometheus.CounterOpts) *MeasurementCounterCollector {
	fqName := prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name)

	collector := MeasurementCounterCollector{
		fqName: fqName,
		opts:   opts,
		desc: prometheus.NewDesc(
			fqName,
			opts.Help,
			measurementCollectorVariableLabels,
			nil,
		),
	}
	collector.values.data = map[string]*measurementResult{}
	collector.metricsMap = prometheus.NewMetricVec(
		collector.desc,
		collector.newMetric,
	)

	return &collector
}

// Describe implements prometheus.Collector's Describe
func (c *MeasurementCounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect implements prometheus.Collector's Collect
func (c *MeasurementCounterCollector) Collect(ch chan<- prometheus.Metric) {
	c.values.RLock()
	defer c.values.RUnlock()

	for sKey, value := range c.values.data {
		lvs := strings.Split(sKey, keySeparator)

		ch <- prometheus.NewMetricWithTimestamp(value.timestamp,
			prometheus.MustNewConstMetric(
				c.desc,
				prometheus.CounterValue,
				value.value,
				lvs...,
			),
		)
	}
}

// Set sets the specified value for provided labelValues at a specified timestamp.
// value must be higher than 0. Otherwise, en error will be logged.
//
// This function is thread-safe.
func (c *MeasurementCounterCollector) Set(timestamp time.Time, value float64, labelValues ...string) {
	if value < 0 {
		log.Println("[ERROR] counters cannot decrease in its value. ignoring.", c.fqName)
		return
	}

	c.values.Lock()
	defer c.values.Unlock()

	lvs := strings.Join(labelValues, keySeparator)
	c.values.data[lvs] = &measurementResult{
		value:     value,
		timestamp: timestamp,
	}
}

// MeasurementGaugeCollector is a struct which takes care of collecting
// prometheus.Gauge metrics Whenever prometheus.Collect is called, all values of
// MeasurementGaugeCollector.values are flushed and sent to the collector
// channel.
type MeasurementGaugeCollector struct {
	metricsMap *prometheus.MetricVec
	desc       *prometheus.Desc
	values     struct {
		sync.RWMutex
		data map[string]*measurementResult // Contains latest value of meters.MeasurementResult
	}
	fqName string
	opts   prometheus.GaugeOpts
}

func NewMeasurementGaugeCollector(opts prometheus.GaugeOpts) *MeasurementGaugeCollector {
	fqName := prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name)

	collector := MeasurementGaugeCollector{
		fqName: fqName,
		opts:   opts,
		desc: prometheus.NewDesc(
			fqName,
			opts.Help,
			measurementCollectorVariableLabels,
			nil,
		),
	}
	collector.values.data = map[string]*measurementResult{}
	collector.metricsMap = prometheus.NewMetricVec(
		collector.desc,
		collector.newMetric,
	)

	return &collector
}

// Describe implements prometheus.Collector's Describe
func (c *MeasurementGaugeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect implements prometheus.Collector's Collect
func (c *MeasurementGaugeCollector) Collect(ch chan<- prometheus.Metric) {
	c.values.RLock()
	defer c.values.RUnlock()

	for sKey, value := range c.values.data {
		lvs := strings.Split(sKey, keySeparator)

		ch <- prometheus.NewMetricWithTimestamp(value.timestamp,
			prometheus.MustNewConstMetric(
				c.desc,
				prometheus.GaugeValue,
				value.value,
				lvs...,
			),
		)
	}
}

// Set sets the specified value for provided labelValues at a specified timestamp.
//
// This function is thread-safe.
func (c *MeasurementGaugeCollector) Set(timestamp time.Time, value float64, labelValues ...string) {
	c.values.Lock()
	defer c.values.Unlock()

	lvs := strings.Join(labelValues, keySeparator)
	c.values.data[lvs] = &measurementResult{
		value:     value,
		timestamp: timestamp,
	}
}

type measurementResult struct {
	timestamp time.Time
	value     float64
}

// Separator used for concatenating label values.
const keySeparator = ";"

// Copied from prometheus/labels.go for consistency purposes
var errInconsistentCardinality = errors.New("inconsistent label cardinality")

// Labels for every measurement prometheus.CounterVec/ unit as e.g. C, F, hertz
var measurementCollectorVariableLabels = []string{"device_name", "serial_number", "unit"}

type MetricFactory interface {
	newMetric(lvs ...string)
}

func (c *MeasurementCounterCollector) newMetric(lvs ...string) prometheus.Metric {
	if len(lvs) != len(measurementCollectorVariableLabels) {
		log.Fatalf(
			"%s: %q has %d variable labels named %q but %d values %q were provided\n",
			errInconsistentCardinality, c.fqName,
			len(measurementCollectorVariableLabels), measurementCollectorVariableLabels,
			len(lvs), lvs,
		)
	}

	return prometheus.NewCounter(c.opts)
}

func (c *MeasurementGaugeCollector) newMetric(lvs ...string) prometheus.Metric {
	if len(lvs) != len(measurementCollectorVariableLabels) {
		log.Fatalf(
			"%s: %q has %d variable labels named %q but %d values %q were provided\n",
			errInconsistentCardinality, c.fqName,
			len(measurementCollectorVariableLabels), measurementCollectorVariableLabels,
			len(lvs), lvs,
		)
	}

	return prometheus.NewGauge(c.opts)
}
