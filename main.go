package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Serializator/magento2-prometheus-exporter-golang/collector"
	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// TODO: support command-line flags for non-default values (e.g. "magento.url" and "magento.bearer" in the YAML configuration)

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	prometheus.MustRegister(collector.NewEnvironmentInfoCollector(
		*httpClient,
		*config,
	))

	prometheus.MustRegister(collector.NewOrdersCollector(
		*httpClient,
		*config,
	))

	http.Handle("/metrics", promhttp.Handler())

	// TODO: make the host to which the exporter will be bound configurable
	// TODO: make the port to which the exporter will be bound configurable
	log.Println("Listening on port :9101")
	log.Fatal(http.ListenAndServe(":9101", nil))
}
