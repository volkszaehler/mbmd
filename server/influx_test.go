package server

import (
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

func TestInfluxBufferOptions(t *testing.T) {
	tc := []struct {
		buffer     time.Duration
		batches    uint
		maxRetries uint
		retryTime  uint
	}{
		// client defaults: 50,000 / 5,000 = 10 batches, 5 retries within 3 minutes
		{0, 10, 5, 180_000},
		// one batch per 1s flush interval, retried indefinitely
		{time.Hour, 3600, math.MaxUint32, math.MaxUint32},
		{10 * time.Minute, 600, math.MaxUint32, math.MaxUint32},
		// shorter than one flush interval still keeps one batch
		{500 * time.Millisecond, 1, math.MaxUint32, math.MaxUint32},
	}

	for _, tc := range tc {
		t.Run(tc.buffer.String(), func(t *testing.T) {
			influx := NewInfluxClient("http://localhost:8086", "db", "mbmd", "org", "", "", "", tc.buffer)
			defer influx.client.Close()

			opts := influx.client.Options().WriteOptions()

			if got := opts.RetryBufferLimit() / opts.BatchSize(); got != tc.batches {
				t.Errorf("retry buffer holds %d batches, want %d", got, tc.batches)
			}
			if got := opts.MaxRetries(); got != tc.maxRetries {
				t.Errorf("max retries %d, want %d", got, tc.maxRetries)
			}
			if got := opts.MaxRetryTime(); got != tc.retryTime {
				t.Errorf("max retry time %d, want %d", got, tc.retryTime)
			}
		})
	}
}

// fakeInflux counts written lines and fails all writes while down is set
type fakeInflux struct {
	down  atomic.Bool
	mu    sync.Mutex
	lines int
}

func (f *fakeInflux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.down.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	body, _ := io.ReadAll(r.Body)
	f.mu.Lock()
	f.lines += len(strings.Split(strings.TrimSpace(string(body)), "\n"))
	f.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (f *fakeInflux) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lines
}

// simulateOutage writes points during a 20s database outage, which is longer
// than the ~10s the client's default retry buffer covers. It keeps writing
// after the database is back until all points arrived or wait has passed,
// and returns the number of points sent and received.
func simulateOutage(t *testing.T, buffer, wait time.Duration) (sent, received int) {
	fake := &fakeInflux{}
	srv := httptest.NewServer(fake)
	defer srv.Close()

	influx := NewInfluxClient(srv.URL, "db", "mbmd", "org", "", "", "", buffer)

	in := make(chan QuerySnip)
	done := make(chan struct{})
	go func() {
		influx.Run(in)
		close(done)
	}()

	send := func(value float64) {
		in <- QuerySnip{
			Device:            "test",
			MeasurementResult: meters.MeasurementResult{Measurement: meters.Power, Value: value},
		}
		sent++
	}

	fake.down.Store(true)

	const points = 100
	for i := range points {
		send(float64(i))
		time.Sleep(200 * time.Millisecond)
	}

	fake.down.Store(false)

	// keep data flowing, the client only retries when it writes new batches
	deadline := time.Now().Add(wait)
	for fake.count() < sent && time.Now().Before(deadline) {
		send(-1)
		time.Sleep(time.Second)
	}

	close(in)
	<-done

	return sent, fake.count()
}

// TestInfluxOutage verifies that points written while the database is down
// are delivered once it is back.
func TestInfluxOutage(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	t.Parallel()

	// the wait covers the client's exponential backoff of up to 125s between retries
	if sent, received := simulateOutage(t, time.Hour, 3*time.Minute); received != sent {
		t.Errorf("received %d of %d points", received, sent)
	}
}

// TestInfluxOutageDefault verifies that without a buffer the client keeps its
// default behaviour and drops data during a longer outage.
func TestInfluxOutageDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	t.Parallel()

	if sent, received := simulateOutage(t, 0, time.Minute); received >= sent {
		t.Errorf("received %d of %d points, expected data loss without buffer", received, sent)
	}
}
