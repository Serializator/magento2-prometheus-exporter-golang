package collector

import (
	"encoding/json"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"strconv"
)

type invoicesCollector struct {
	// Dependencies for this collector are defined below
	client magento.Client

	// Descriptors for this collector are defined below
	total *prometheus.GaugeVec
}

func NewInvoicesCollector(client magento.Client) *invoicesCollector {
	return &invoicesCollector{
		client: client,

		total: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "magento",
			Subsystem: "invoices",
			Name:      "total",
			Help:      "Total amount of invoices",
		}, []string{"store_id", "state"}),
	}
}

func (collector *invoicesCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, descs)
}

func (collector *invoicesCollector) Collect(metrics chan<- prometheus.Metric) {
	invoicesResponse, err := collector.fetchAndDecodeInvoices()
	if err != nil {
		descs := make(chan *prometheus.Desc, 1)
		collector.total.Describe(descs)
		metrics <- prometheus.NewInvalidMetric(<-descs, err)
		return
	}

	collector.total.Reset()

	for _, invoiceMetricAggregation := range invoicesResponse.Total {
		counter, err := collector.total.GetMetricWithLabelValues(strconv.Itoa(invoiceMetricAggregation.StoreId), invoiceMetricAggregation.State)
		if err != nil {
			metrics <- prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Set(float64(invoiceMetricAggregation.Count))
	}

	collector.total.Collect(metrics)
}

func (collector *invoicesCollector) fetchAndDecodeInvoices() (invoicesResponse, error) {
	invoicesResponse := &invoicesResponse{}

	response, err := collector.client.Get("/metrics/invoices")
	if err != nil {
		return *invoicesResponse, err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(invoicesResponse); err != nil {
		return *invoicesResponse, err
	}

	return *invoicesResponse, nil
}
