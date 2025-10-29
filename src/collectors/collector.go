package collectors

import (
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const exporterNamespace = "ds"

// InternalCollector describes the interface of the internal collector that collects data.
type InternalCollector interface {
	Get(chan<- prometheus.Metric) error
}

// DSCollector implements prometheus.Collector interface.
type DSCollector struct {
	collectors         map[string]InternalCollector
	scrapeDurationDesc *prometheus.Desc
	scrapeSuccessDesc  *prometheus.Desc
}

func NewDSCollector() *DSCollector {
	collector := DSCollector{
		collectors: make(map[string]InternalCollector),
		scrapeDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(exporterNamespace, "scrape", "duration_seconds"),
			"Duration of a collector scrape",
			[]string{"collector"},
			nil,
		),
		scrapeSuccessDesc: prometheus.NewDesc(
			prometheus.BuildFQName(exporterNamespace, "scrape", "success"),
			"Whether a collector succeeded",
			[]string{"collector"},
			nil,
		),
	}
	return &collector
}

// Describe DSCollector metrics.
func (c *DSCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- c.scrapeDurationDesc
	channel <- c.scrapeSuccessDesc
}

// Collect initiates the receipt of metrics from all collectors.
func (c *DSCollector) Collect(channel chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(len(c.collectors))
	for collector := range c.collectors {
		go func() {
			c.scrape(collector, channel)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (c *DSCollector) Register(name string, collector InternalCollector) {
	c.collectors[name] = collector
}

func (c *DSCollector) scrape(collector string, channel chan<- prometheus.Metric) {
	start_time := time.Now()

	err := c.collectors[collector].Get(channel)

	elapsed := time.Since(start_time)
	channel <- prometheus.MustNewConstMetric(
		c.scrapeDurationDesc,
		prometheus.GaugeValue,
		float64(elapsed.Seconds()),
		collector)

	if err != nil {
		slog.Error("Colletor failed", "collector", collector, "err", err)
		channel <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, 0, collector)
	} else {
		channel <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, 1, collector)
	}

}
