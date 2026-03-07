/*
The collectors package provides various structures implementing the prometheus.Collector interface
*/
package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"

	expldap "389-ds-exporter/internal/ldap"
)

const ldapTimestampLayout = "20060102150405Z"

type Parser func(string) (float64, error)

func ParseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func ParseTimestamp(value string) (float64, error) {
	parsedTime, err := time.Parse(ldapTimestampLayout, value)
	if err != nil {
		return 0, err
	}
	return float64(parsedTime.Unix()), nil
}

// LdapMetric implements a container for storing
// the necessary information about an attribute used in metrics.
type LdapMetric struct {
	MetricName string
	LdapName   string
	Help       string
	Type       prometheus.ValueType
	IsInfo     bool
	Parser     Parser
	//Labels    prometheus.Labels
}

// LdapEntryCollector collects 389-ds metrics.
type LdapEntryCollector struct {
	connectionPool *expldap.Pool
	baseDn         string
	metrics        []LdapMetric
	descriptors    map[string]*prometheus.Desc
	poolGetTimeout time.Duration
}

// NewLdapEntryCollector function create new LdapEntryCollector instance based on provided parameters.
func NewLdapEntryCollector(
	subsystem string,
	connectionPool *expldap.Pool,
	entryBaseDn string,
	metrics []LdapMetric,
	labels prometheus.Labels,
	poolGetTimeout time.Duration,
) *LdapEntryCollector {
	metricsDescriptors := make(map[string]*prometheus.Desc)

	for idx, metric := range metrics {
		if metric.MetricName == "" {
			panic("LdapMetric.MetricName cannot be empty")
		}

		if metric.Parser == nil {
			metrics[idx].Parser = ParseFloat
		}

		var labelNames []string

		if metric.IsInfo {
			labelNames = []string{"value"}
		}

		metricsDescriptors[metric.MetricName] = prometheus.NewDesc(
			prometheus.BuildFQName(exporterNamespace, subsystem, metric.MetricName),
			metric.Help,
			labelNames,
			labels,
		)
	}

	return &LdapEntryCollector{
		connectionPool: connectionPool,
		baseDn:         entryBaseDn,
		metrics:        metrics,
		descriptors:    metricsDescriptors,
		poolGetTimeout: poolGetTimeout,
	}
}

// Get function fetches metrics from LDAP and sends them to the provided channel.
func (c *LdapEntryCollector) Get(channel chan<- prometheus.Metric) error {

	ldapEntries, err := c.getLdapEntryAttributes()
	if err != nil {
		return fmt.Errorf("error getting attrs from LDAP: %w", err)
	}

	var result error

	for _, metric := range c.metrics {

		raw, ok := ldapEntries[metric.LdapName]
		if !ok {
			slog.Debug("Attribute not found in LDAP response", "attr", metric.LdapName)
			continue
		}

		desc := c.descriptors[metric.MetricName]

		if metric.IsInfo {
			channel <- prometheus.MustNewConstMetric(desc, metric.Type, 1, raw)
			continue
		}

		value, err := metric.Parser(raw)
		if err != nil {
			slog.Debug(
				"Metric parsing failed",
				"attr", metric.LdapName,
				"value", raw,
				"err", err,
			)
			result = fmt.Errorf("metric parsing error: %w", err)
			continue
		}

		channel <- prometheus.MustNewConstMetric(desc, metric.Type, value)
	}

	return result
}

// getLdapEntryAttributes returns the record attributes specified in the LdapEntryCollector from the ldap.
func (c *LdapEntryCollector) getLdapEntryAttributes() (map[string]string, error) {
	attributeList := make([]string, 0, len(c.metrics))

	for _, monitoredAttr := range c.metrics {
		attributeList = append(attributeList, monitoredAttr.LdapName)
	}

	searchAttributesRequest := ldap.NewSearchRequest(
		c.baseDn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectclass=*)",
		attributeList,
		nil,
	)

	ctx, cancel := context.WithTimeout(context.Background(), c.poolGetTimeout)
	defer cancel()
	conn, err := c.connectionPool.Conn(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get connection from pool: %w", err)
	}
	defer conn.Close()

	searchResult, err := conn.Search(searchAttributesRequest)
	if err != nil {
		return nil, fmt.Errorf(
			"LDAP Search request (dn='%v', attrs='%v') failed with error: %w",
			searchAttributesRequest.BaseDN,
			searchAttributesRequest.Attributes,
			err,
		)
	}

	attrs := make(map[string]string)

	if len(searchResult.Entries) < 1 {
		slog.Warn("LDAP request returned no entries. The configuration may be incorrect or the user may not have permissions",
			"req_dn", searchAttributesRequest.BaseDN,
			"req_attrs", searchAttributesRequest.Attributes)
		return attrs, nil
	}

	for _, attr := range searchResult.Entries[0].Attributes {
		if len(attr.Values) > 1 {
			slog.Debug("Attribute has multiple values, using first", "attr", attr.Name)
		}
		attrs[attr.Name] = attr.Values[0]
	}

	return attrs, nil
}
