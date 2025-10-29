package collectors

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/connections"
)

// PoolCollector collects internal LDAP-pool metrics.
type PoolCollector struct {
	connectionPool   *connections.LDAPPool
	descOpen         *prometheus.Desc
	descClosedIdle   *prometheus.Desc
	descClosedLife   *prometheus.Desc
	descWaitCount    *prometheus.Desc
	descWaitDuration *prometheus.Desc
	mutex            sync.Mutex
}

// NewPoolCollector function create new PoolCollector instance based on provided parameters.
func NewPoolCollector(
	subsystem string,
	connectionPool *connections.LDAPPool,
	labels prometheus.Labels,
) *PoolCollector {
	pool := &PoolCollector{connectionPool: connectionPool}

	pool.descOpen = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, subsystem, "open"),
		"Number of open connections in the pool.",
		nil,
		labels,
	)
	pool.descClosedIdle = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, subsystem, "closed_idletime"),
		"Number of connections closed after idle timeout.",
		nil,
		labels,
	)

	pool.descClosedLife = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, subsystem, "closed_lifetime"),
		"Number of connections closed after lifetime expired.",
		nil,
		labels,
	)

	pool.descWaitCount = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, subsystem, "wait_count"),
		"The number of times clients waited for a connection to appear in the pool.",
		nil,
		labels,
	)

	pool.descWaitDuration = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, subsystem, "wait_duration"),
		"Total time spent waiting for a connection",
		nil,
		labels,
	)

	return pool
}

// Get function fetches metrics from LDAP and sends them to the provided channel.
func (c *PoolCollector) Get(channel chan<- prometheus.Metric) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stat := c.connectionPool.Stat()

	channel <- prometheus.MustNewConstMetric(c.descOpen, prometheus.GaugeValue, float64(stat.Open))
	channel <- prometheus.MustNewConstMetric(c.descClosedIdle, prometheus.CounterValue, float64(stat.ClosedIdleTime))
	channel <- prometheus.MustNewConstMetric(c.descClosedLife, prometheus.CounterValue, float64(stat.ClosedLifeTime))
	channel <- prometheus.MustNewConstMetric(c.descWaitCount, prometheus.CounterValue, float64(stat.WaitCount))
	channel <- prometheus.MustNewConstMetric(
		c.descWaitDuration,
		prometheus.CounterValue,
		float64(time.Duration(stat.WaitDuration).Milliseconds()),
	)
	return nil
}
