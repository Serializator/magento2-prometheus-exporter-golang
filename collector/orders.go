package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
)

type ordersCollector struct {
	// Dependencies for this collector are defined below
	http   http.Client
	config config.Config

	// Descriptors for this collector are defined below
	total *prometheus.GaugeVec
}

func NewOrdersCollector(http http.Client, config config.Config) *ordersCollector {
	return &ordersCollector{
		http:   http,
		config: config,

		total: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "magento",
			Subsystem: "orders",
			Name:      "total",
			Help:      "Total amount of orders",
		}, []string{"state", "status"}),
	}
}

func (collector *ordersCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, descs)
}

func (collector *ordersCollector) Collect(metrics chan<- prometheus.Metric) {
	ordersResponse, err := collector.fetchAndDecodeOrders()
	if err != nil {
		// TODO: use "prometheus.NewInvalidMetric"
		return
	}

	collector.total.Reset()

	for _, order := range ordersResponse.Items {
		counter, err := collector.total.GetMetricWithLabelValues(order.State, order.Status)
		if err != nil {
			prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Inc()
	}

	collector.total.Collect(metrics)
}

func (collector *ordersCollector) fetchAndDecodeOrders() (ordersResponse, error) {
	// TODO: refactor HTTP requests such that the Magento URL and authorization code can be re-used

	ordersResponse := &ordersResponse{}

	queryString := []string{
		"searchCriteria[filter_groups][0][filters][0][field]=entity_id",
		"searchCriteria[filter_groups][0][filters][0][value]=0",
		"searchCriteria[filter_groups][0][filters][0][condition_type]=gt",
		"fields=items[status,state]",
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/V1/orders?%s",
		collector.config.Magento.Url,
		strings.Join(queryString, "&"),
	), nil)
	if err != nil {
		return *ordersResponse, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", collector.config.Magento.Bearer))
	response, err := collector.http.Do(request)
	if err != nil {
		return *ordersResponse, err
	}

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(ordersResponse); err != nil {
		return *ordersResponse, err
	}

	return *ordersResponse, nil
}
