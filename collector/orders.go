package collector

import (
	"encoding/json"
	"fmt"
	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
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
		}, []string{"store_id", "state", "status", "payment_method"}),
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

	for _, order := range ordersResponse.Total {
		counter, err := collector.total.GetMetricWithLabelValues(strconv.Itoa(order.StoreId), order.State, order.Status, order.PaymentMethod)
		if err != nil {
			prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Set(float64(order.Count))
	}

	collector.total.Collect(metrics)
}

func (collector *ordersCollector) fetchAndDecodeOrders() (*ordersResponse, error) {
	// TODO: refactor HTTP requests such that the Magento URL and authorization code can be re-used

	ordersResponse := &ordersResponse{}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/V1/metrics/orders", collector.config.Magento.Url), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", collector.config.Magento.Bearer))
	response, err := collector.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(ordersResponse); err != nil {
		return nil, err
	}

	return ordersResponse, nil
}
