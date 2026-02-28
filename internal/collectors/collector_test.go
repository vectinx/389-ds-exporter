package collectors

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

type successCollector struct{}

func (s *successCollector) Get(ch chan<- prometheus.Metric) error {
	return nil
}

type errorCollector struct{}

func (e *errorCollector) Get(ch chan<- prometheus.Metric) error {
	return errors.New("something went wrong")
}

func disableLogs(t *testing.T) {
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	t.Cleanup(func() {
		slog.SetDefault(original)
	})
}

func TestNewDSCollectorInit(t *testing.T) {
	c := NewDSCollector()

	require.NotNil(t, c, "collector should not be nil")
	require.NotNil(t, c.collectors, "c.collectors should not be nil")
	require.NotNil(t, c.scrapeDurationDesc, "c.scrapeDurationDesc should not be nil")
	require.NotNil(t, c.scrapeSuccessDesc, "c.scrapeSuccessDesc should not be nil")
}

func TestDSCollectorRegister(t *testing.T) {
	c := NewDSCollector()
	c.Register("test", &successCollector{})

	require.Contains(t, c.collectors, "test", "Collector 'test' must be registered")
}

func TestDSCollectorCollectSuccess(t *testing.T) {
	c := NewDSCollector()
	c.Register("test", &successCollector{})

	expected := `
# HELP ds_exporter_scrape_success Whether a collector succeeded
# TYPE ds_exporter_scrape_success gauge
ds_exporter_scrape_success{collector="test"} 1
`

	err := testutil.CollectAndCompare(
		c,
		strings.NewReader(expected),
		"ds_exporter_scrape_success",
	)
	require.NoError(t, err, "Collect should not fail")

	count := testutil.CollectAndCount(c, "ds_exporter_scrape_duration_seconds")
	require.Equal(t, 1, count, "There should be one duration metric")
}

func TestDSCollectorCollectError(t *testing.T) {
	// Turn off logging so as not to clog up the test output.
	disableLogs(t)

	c := NewDSCollector()
	c.Register("test", &errorCollector{})

	expected := `
# HELP ds_exporter_scrape_success Whether a collector succeeded
# TYPE ds_exporter_scrape_success gauge
ds_exporter_scrape_success{collector="test"} 0
`

	err := testutil.CollectAndCompare(
		c,
		strings.NewReader(expected),
		"ds_exporter_scrape_success",
	)
	require.NoError(t, err, "Collect should not fail")

	count := testutil.CollectAndCount(c, "ds_exporter_scrape_duration_seconds")
	require.Equal(t, 1, count, "There should be one duration metric")
}

func TestDSCollectorCollectMultipleCollectors(t *testing.T) {
	c := NewDSCollector()
	c.Register("one", &successCollector{})
	c.Register("two", &errorCollector{})

	expected := `
# HELP ds_exporter_scrape_success Whether a collector succeeded
# TYPE ds_exporter_scrape_success gauge
ds_exporter_scrape_success{collector="one"} 1
ds_exporter_scrape_success{collector="two"} 0
`

	err := testutil.CollectAndCompare(
		c,
		strings.NewReader(expected),
		"ds_exporter_scrape_success",
	)
	require.NoError(t, err, "Collect should not fail")

	count := testutil.CollectAndCount(c, "ds_exporter_scrape_duration_seconds")
	require.Equal(t, 2, count, "There should be one duration metric")
}

func TestDSCollectorCollectNoCollectors(t *testing.T) {
	c := NewDSCollector()

	ch := make(chan prometheus.Metric)
	done := make(chan struct{})

	go func() {
		c.Collect(ch)
		close(ch)
		close(done)
	}()

	for range ch {
	}

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond, "Collect blocked with no collectors")
}

func TestDSCollector_CollectNoCollectors(t *testing.T) {
	c := NewDSCollector()

	ch := make(chan prometheus.Metric)
	done := make(chan struct{})

	go func() {
		c.Collect(ch)
		close(ch)
		close(done)
	}()

	for range ch {
	}

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond, "Collect blocked with no collectors")
}
