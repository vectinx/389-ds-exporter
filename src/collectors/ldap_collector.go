/*
The collectors package provides various structures implementing the prometheus.Collector interface
*/
package collectors

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"sync"
	"time"

	"389-ds-exporter/src/backends"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"
)

/*
LdapValueType type defines the format
in which the attribute value is stored in the LDAP and how it will be converted to float64.
*/
type LdapValueType int

const dateTimeLayout string = "20060102150405Z"

const (
	_ LdapValueType = iota
	// NumericValue type corresponds a simple numeric value.
	NumericValue
	// Iso8601CompactString type corresponds a string containing the date and time in the 'YYYYMMDDThhmmssZ' format.
	Iso8601CompactString
)

// LdapMonitoredAttribute implements a container for storing the necessary information about an attribute used in metrics.
type LdapMonitoredAttribute struct {
	LdapName string
	LdapType LdapValueType
	Help     string
	Type     prometheus.ValueType
	Labels   prometheus.Labels
}

// LdapCollector collects 389-ds metrics. It implements prometheus.Collector interface.
type LdapEntryCollector struct {
	connectionPool *backends.LdapConnectionPool
	namespace      string
	baseDn         string
	attributes     map[string]LdapMonitoredAttribute
	descriptors    map[string]*prometheus.Desc
	mutex          sync.Mutex
}

// NewLdapEntryCollector function create new LdapEntryCollector instance based on provided parameteres.
func NewLdapEntryCollector(
	namespace string,
	connectionPool *backends.LdapConnectionPool,
	entryBaseDn string,
	attributes map[string]LdapMonitoredAttribute,
	labels prometheus.Labels,
) *LdapEntryCollector {
	metricsDescriptors := make(map[string]*prometheus.Desc)

	for key, val := range attributes {
		metricsDescriptors[key] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", key),
			val.Help,
			nil,
			labels,
		)
	}

	return &LdapEntryCollector{
		connectionPool: connectionPool,
		namespace:      namespace,
		baseDn:         entryBaseDn,
		attributes:     attributes,
		descriptors:    metricsDescriptors,
	}
}

func (c *LdapEntryCollector) getLdapEntryAttributes() (map[string][]string, error) {
	ldapConnection, err := c.connectionPool.Get(5 * time.Second)
	defer c.connectionPool.Put(ldapConnection)

	if err != nil {
		return nil, fmt.Errorf("error getting LDAP connection from pool: %w", err)
	}

	if ldapConnection == nil {
		log.Print("А соединение почему-то пустое ...")
	}
	var attributeList []string

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

	searchResult, err := ldapConnection.Conn.Search(searchAttributesRequest)
	if err != nil {
		return nil, fmt.Errorf("LDAP Search request failed with error: %w", err)
	}

	returnValue := make(map[string][]string)

	for _, attr := range searchResult.Entries[0].Attributes {
		if !slices.Contains(attributeList, attr.Name) {
			continue
		}

		returnValue[attr.Name] = attr.Values
	}

	return returnValue, nil
}

// Describe function sends the super-set of all possible descriptors of LDAP metrics
// to the provided channel.
func (c *LdapEntryCollector) Describe(channel chan<- *prometheus.Desc) {
	for _, descriptor := range c.descriptors {
		channel <- descriptor
	}
}

// Collect function fetches metrics from LDAP and sends them to the provided channel.
func (c *LdapEntryCollector) Collect(channel chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ldapEntries, err := c.getLdapEntryAttributes()
	if err != nil {
		log.Printf("Error getting values from ldap: %s", err)
		return
	}

	for key, value := range c.attributes {
		attributeValues, ok := ldapEntries[value.LdapName]
		if !ok {
			log.Printf("Attribute %v was not found in LDAP response. ", value.LdapName)
			continue
		}

		if len(attributeValues) > 1 {
			log.Printf("Attribute %s has more than one value, the first one will be used", key)
		}

		var converted float64

		if value.LdapType == Iso8601CompactString {
			parsedTime, err := time.Parse(dateTimeLayout, attributeValues[0])
			if err != nil {
				log.Printf("Error parsing date: %s", err)

				continue
			}

			converted = float64(parsedTime.Unix())
		} else {
			converted, err = strconv.ParseFloat(ldapEntries[value.LdapName][0], 64)
			if err != nil {
				log.Printf("Unable to convert attribute \"%s\" value \"%s\" to type float64", key, ldapEntries[value.LdapName])

				continue
			}
		}

		channel <- prometheus.MustNewConstMetric(c.descriptors[key],
			c.attributes[key].Type, converted)
	}
}
