package collector

import (
	"fmt"
	"github.com/Serializator/magento2-prometheus-exporter-golang/magento"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Serializator/magento2-prometheus-exporter-golang/config"
)

func TestEnvironmentInfo(t *testing.T) {
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"version":"2.4.3","edition":"Commerce","mode":"production"}`)
	}))

	defer testHttpServer.Close()

	testHttpServerUrl, err := url.Parse(testHttpServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse URL of the HTTP test server: %s", err)
	}

	environmentInfo := NewEnvironmentInfoCollector(*magento.NewClient(
		testHttpServerUrl.String(),
		magento.NewBearerAuthenticator(""),
		http.DefaultClient,
	), config.Config{
		Magento: struct {
			Url    string
			Bearer string
		}{
			Url:    testHttpServerUrl.String(),
			Bearer: "",
		},
	})

	environmentInfoResponse, err := environmentInfo.fetchAndDecodeEnvironmentInfo()
	if err != nil {
		t.Fatalf("Failed to fetch or decode environment info: %s", err)
	}

	if environmentInfoResponse.Version != "2.4.3" {
		t.Errorf("Invalid Magento version response")
	}

	if environmentInfoResponse.Edition != "Commerce" {
		t.Errorf("Invalid Magento edition response")
	}

	if environmentInfoResponse.Mode != "production" {
		t.Errorf("Invalid Magento mode response")
	}
}
