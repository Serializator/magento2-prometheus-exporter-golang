package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
)

type environmentInfoCollector struct {
	// Dependencies for this collector are defined below
	http   http.Client
	config config.Config

	// Descriptors for this collector are defined below
	up   prometheus.Gauge
	info *prometheus.Desc
}

func NewEnvironmentInfoCollector(http http.Client, config config.Config) *environmentInfoCollector {
	return &environmentInfoCollector{
		http:   http,
		config: config,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("magento", "environment", "up"),
			Help: "Was the last scrape of the Magento environment successful.",
		}),

		info: prometheus.NewDesc(
			prometheus.BuildFQName("magento", "environment", "info"),
			"Information about the Magento environment from which the exporter is pulling its metrics.",
			[]string{"version", "edition", "mode"}, nil,
		),
	}
}

func (collector *environmentInfoCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, descs)
}

func (collector *environmentInfoCollector) Collect(metrics chan<- prometheus.Metric) {
	defer func() {
		metrics <- collector.up
	}()

	environmentInfo, err := collector.fetchAndDecodeEnvironmentInfo()
	if err != nil {
		collector.up.Set(0)
		return
	}

	collector.up.Set(1)

	metrics <- prometheus.MustNewConstMetric(
		collector.info,
		prometheus.GaugeValue, 1,
		environmentInfo.Version, environmentInfo.Edition, environmentInfo.Mode,
	)
}

func (collector *environmentInfoCollector) fetchAndDecodeEnvironmentInfo() (environmentInfoResponse, error) {
	environmentInfoResponse := &environmentInfoResponse{}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/V1/metrics/info", collector.config.Magento.Url), nil)
	if err != nil {
		return *environmentInfoResponse, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", collector.config.Magento.Bearer))
	response, err := collector.http.Do(request)
	if err != nil {
		return *environmentInfoResponse, err
	}

	defer response.Body.Close()
	if err = json.NewDecoder(response.Body).Decode(environmentInfoResponse); err != nil {
		return *environmentInfoResponse, err
	}

	return *environmentInfoResponse, nil
}
