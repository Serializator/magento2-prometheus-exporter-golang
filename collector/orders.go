package collector

import (
	"encoding/json"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"strconv"
)

type ordersCollector struct {
	// Dependencies for this collector are defined below
	client magento.Client

	// Descriptors for this collector are defined below
	total *prometheus.GaugeVec
}

func NewOrdersCollector(client magento.Client) *ordersCollector {
	return &ordersCollector{
		client: client,

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
		descs := make(chan *prometheus.Desc, 1)
		collector.total.Describe(descs)
		metrics <- prometheus.NewInvalidMetric(<-descs, err)
		return
	}

	collector.total.Reset()

	for _, order := range ordersResponse.Total {
		counter, err := collector.total.GetMetricWithLabelValues(strconv.Itoa(order.StoreId), order.State, order.Status, order.PaymentMethod)
		if err != nil {
			metrics <- prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Set(float64(order.Count))
	}

	collector.total.Collect(metrics)
}

func (collector *ordersCollector) fetchAndDecodeOrders() (*ordersResponse, error) {
	ordersResponse := &ordersResponse{}

	response, err := collector.client.Get("/metrics/orders")
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(ordersResponse); err != nil {
		return nil, err
	}

	return ordersResponse, nil
}
