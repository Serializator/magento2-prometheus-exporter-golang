package collector

import (
	"encoding/json"
	"fmt"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"io"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
		prometheus.NewInvalidMetric(<-descs, err)
		return
	}

	collector.total.Reset()

	for _, invoice := range invoicesResponse.Items {
		state, err := invoice.State.String()
		if err != nil {
			// TODO: use "prometheus.NewInvalidMetric"
			continue
		}

		counter, err := collector.total.GetMetricWithLabelValues(strconv.Itoa(invoice.StoreId), state)
		if err != nil {
			prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Inc()
	}

	collector.total.Collect(metrics)
}

func (collector *invoicesCollector) fetchAndDecodeInvoices() (invoicesResponse, error) {
	invoicesResponse := &invoicesResponse{}

	queryString := []string{
		"searchCriteria[filter_groups][0][filters][0][field]=entity_id",
		"searchCriteria[filter_groups][0][filters][0][value]=0",
		"searchCriteria[filter_groups][0][filters][0][condition_type]=gt",
		"fields=items[store_id,state]",
	}
	response, err := collector.client.Get(fmt.Sprintf("/invoices?%s", strings.Join(queryString, "&")))
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
