/*
The collectors package provides various structures implementing the prometheus.Collector interface
*/
package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"

	expldap "389-ds-exporter/internal/ldap"
)

/*
LdapAttrValueType type defines the format
in which the attribute value is stored in the LDAP and how it will be converted to float64.
*/
type LdapAttrValueType int

const dateTimeLayout string = "20060102150405Z"

const (
	_ LdapAttrValueType = iota
	// NumericValue type corresponds a simple numeric value.
	NumericValue
	// Iso8601CompactString type corresponds a string containing the date and time in the 'YYYYMMDDThhmmssZ' format.
	Iso8601CompactString
	// StringLabel type corresponds ещ a string value that should be presented as a variable label.
	StringLabel
)

// LdapMonitoredAttribute implements a container for storing
// the necessary information about an attribute used in metrics.
type LdapMonitoredAttribute struct {
	LdapName string
	LdapType LdapAttrValueType
	Help     string
	Type     prometheus.ValueType
	Labels   prometheus.Labels
}

// LdapEntryCollector collects 389-ds metrics.
type LdapEntryCollector struct {
	connectionPool *expldap.Pool
	baseDn         string
	attributes     map[string]LdapMonitoredAttribute
	descriptors    map[string]*prometheus.Desc
	mutex          sync.Mutex
	poolGetTimeout time.Duration
}

// NewLdapEntryCollector function create new LdapEntryCollector instance based on provided parameters.
func NewLdapEntryCollector(
	subsystem string,
	connectionPool *expldap.Pool,
	entryBaseDn string,
	attributes map[string]LdapMonitoredAttribute,
	labels prometheus.Labels,
	poolGetTimeout time.Duration,
) *LdapEntryCollector {
	metricsDescriptors := make(map[string]*prometheus.Desc)

	for key, val := range attributes {
		if val.LdapType == StringLabel {
			metricsDescriptors[key] = prometheus.NewDesc(
				prometheus.BuildFQName(exporterNamespace, subsystem, key),
				val.Help,
				[]string{key},
				labels,
			)
		} else {
			metricsDescriptors[key] = prometheus.NewDesc(
				prometheus.BuildFQName(exporterNamespace, subsystem, key),
				val.Help,
				nil,
				labels,
			)
		}
	}

	return &LdapEntryCollector{
		connectionPool: connectionPool,
		baseDn:         entryBaseDn,
		attributes:     attributes,
		descriptors:    metricsDescriptors,
		poolGetTimeout: poolGetTimeout,
	}
}

// Get function fetches metrics from LDAP and sends them to the provided channel.
func (c *LdapEntryCollector) Get(channel chan<- prometheus.Metric) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ldapEntries, err := c.getLdapEntryAttributes()
	if err != nil {
		return fmt.Errorf("error getting attrs from LDAP: %w", err)
	}
	var result error = nil

	for key, value := range c.attributes {

		attributeValues, ok := ldapEntries[value.LdapName]
		var labelValues = []string{}

		if !ok {
			slog.Debug("Attribute was not found in LDAP response. ", "attr_name", value.LdapName)

			continue
		}

		if len(attributeValues) > 1 {
			slog.Debug("Attribute has more than one value, the first one will be used", "attr_name", key)
		}

		var converted float64

		if value.LdapType == Iso8601CompactString {
			parsedTime, err := time.Parse(dateTimeLayout, attributeValues[0])
			if err != nil {
				slog.Debug(
					"Error converting date to type float64",
					"attr_name",
					key,
					"attr_value",
					attributeValues[0],
					"err",
					err,
				)

				result = fmt.Errorf(
					"error converting attribute value to float64: %w",
					err,
				)
				continue
			}

			converted = float64(parsedTime.Unix())
		} else if value.LdapType == StringLabel {
			converted = 1
			labelValues = append(labelValues, attributeValues[0])
		} else {
			converted, err = strconv.ParseFloat(ldapEntries[value.LdapName][0], 64)
			if err != nil {
				slog.Debug(
					"Error converting attribute value to type float64",
					"attr_name",
					key,
					"attr_value",
					ldapEntries[value.LdapName],
				)
				result = fmt.Errorf(
					"error converting attribute value to float64: %w",
					err,
				)
				continue
			}
		}

		channel <- prometheus.MustNewConstMetric(c.descriptors[key],
			c.attributes[key].Type, converted, labelValues...)
	}
	return result
}

// getLdapEntryAttributes returs the record attributes specified in the LdapEntryCollector from the ldap.
func (c *LdapEntryCollector) getLdapEntryAttributes() (map[string][]string, error) {
	attributeList := make([]string, 0, len(c.attributes))

	for _, monitoredAttr := range c.attributes {
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

	returnValue := make(map[string][]string)

	if len(searchResult.Entries) < 1 {
		slog.Warn("LDAP request returned no entries. The configuration may be incorrect or the user may not have permissions",
			"req_dn", searchAttributesRequest.BaseDN,
			"req_attrs", searchAttributesRequest.Attributes)
		return returnValue, nil
	}

	for _, attr := range searchResult.Entries[0].Attributes {
		if !slices.Contains(attributeList, attr.Name) {
			continue
		}

		returnValue[attr.Name] = attr.Values
	}

	return returnValue, nil
}
