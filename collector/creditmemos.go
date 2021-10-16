package collector

import (
	"encoding/json"
	"fmt"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"io"
	"strconv"
	"strings"

	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
)

type creditmemosCollector struct {
	// Dependencies for this collector are defined below
	client magento.Client
	config config.Config

	// Descriptors for this collector are defined below
	total *prometheus.GaugeVec
}

func NewCreditmemosCollector(client magento.Client, config config.Config) *creditmemosCollector {
	return &creditmemosCollector{
		client: client,
		config: config,

		total: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "magento",
			Subsystem: "creditmemos",
			Name:      "total",
			Help:      "Total amount of creditmemos",
		}, []string{"store_id", "state"}),
	}
}

func (collector *creditmemosCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, descs)
}

func (collector *creditmemosCollector) Collect(metrics chan<- prometheus.Metric) {
	creditMemosResponse, err := collector.fetchAndDecodeCreditMemos()
	if err != nil {
		// TODO: use "prometheus.NewInvalidMetric"
		return
	}

	collector.total.Reset()

	for _, creditmemo := range creditMemosResponse.Items {
		state, err := creditmemo.State.String()
		if err != nil {
			// TODO: use "prometheus.NewInvalidMetric"
			continue
		}

		counter, err := collector.total.GetMetricWithLabelValues(strconv.FormatInt(creditmemo.StoreId, 10), state)
		if err != nil {
			prometheus.NewInvalidMetric(counter.Desc(), err)
			continue
		}

		counter.Inc()
	}

	collector.total.Collect(metrics)
}

func (collector *creditmemosCollector) fetchAndDecodeCreditMemos() (creditmemosResponse, error) {
	creditmemosResponse := &creditmemosResponse{}

	queryString := []string{
		"searchCriteria[filter_groups][0][filters][0][field]=entity_id",
		"searchCriteria[filter_groups][0][filters][0][value]=0",
		"searchCriteria[filter_groups][0][filters][0][condition_type]=gt",
		"fields=items[store_id,state]",
	}
	response, err := collector.client.Get(fmt.Sprintf("/creditmemos?%s", strings.Join(queryString, "&")))
	if err != nil {
		return *creditmemosResponse, err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(creditmemosResponse); err != nil {
		return *creditmemosResponse, err
	}

	return *creditmemosResponse, nil
}
