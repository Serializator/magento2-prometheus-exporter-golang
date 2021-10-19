package collector

import (
	"encoding/json"
	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

type environmentInfoCollector struct {
	// Dependencies for this collector are defined below
	client magento.Client
	config config.Config

	// Descriptors for this collector are defined below
	up   prometheus.Gauge
	info *prometheus.Desc
}

func NewEnvironmentInfoCollector(client magento.Client, config config.Config) *environmentInfoCollector {
	return &environmentInfoCollector{
		client: client,
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
		prometheus.NewInvalidMetric(collector.up.Desc(), err)
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

	response, err := collector.client.Get("/metrics/info")
	if err != nil {
		return *environmentInfoResponse, err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
	}(response.Body)

	if err = json.NewDecoder(response.Body).Decode(environmentInfoResponse); err != nil {
		return *environmentInfoResponse, err
	}

	return *environmentInfoResponse, nil
}
