package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
)

type invoicesCollector struct {
	// Dependencies for this collector are defined below
	http   http.Client
	config config.Config

	// Descriptors for this collector are defined below
	total *prometheus.GaugeVec
}

func NewInvoicesCollector(http http.Client, config config.Config) *invoicesCollector {
	return &invoicesCollector{
		http:   http,
		config: config,

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
		// TODO: use "prometheus.NewInvalidMetric"
		return
	}

	collector.total.Reset()

	for _, invoice := range invoicesResponse.Items {
		state, err := invoice.State.String()
		if err != nil {
			// TODO: use "prometheus.NewInvalidMetric"
			continue
		}

		counter, err := collector.total.GetMetricWithLabelValues(strconv.FormatInt(invoice.StoreId, 10), state)
		if err != nil {
			prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Inc()
	}

	collector.total.Collect(metrics)
}

func (collector *invoicesCollector) fetchAndDecodeInvoices() (invoicesResponse, error) {
	// TODO: refactor HTTP requests such that the Magento URL and authorization code can be re-used

	invoicesResponse := &invoicesResponse{}

	queryString := []string{
		"searchCriteria[filter_groups][0][filters][0][field]=entity_id",
		"searchCriteria[filter_groups][0][filters][0][value]=0",
		"searchCriteria[filter_groups][0][filters][0][condition_type]=gt",
		"fields=items[store_id,state]",
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/V1/invoices?%s",
		collector.config.Magento.Url,
		strings.Join(queryString, "&"),
	), nil)
	if err != nil {
		return *invoicesResponse, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", collector.config.Magento.Bearer))
	response, err := collector.http.Do(request)
	if err != nil {
		return *invoicesResponse, err
	}

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(invoicesResponse); err != nil {
		return *invoicesResponse, err
	}

	return *invoicesResponse, nil
}
