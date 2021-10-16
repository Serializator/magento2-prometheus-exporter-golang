package magento

import (
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	url        string
	authenticator Authenticator
}

func NewClient(url string, authenticator Authenticator, httpClient *http.Client) *Client {
	return &Client{httpClient, strings.Trim(url, "/"), authenticator}
}

func (client *Client) Get(path string) (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", client.url, strings.Trim(path, "/")), nil)
	if err != nil {
		return nil, err
	}

	err = client.authenticator.apply(request)
	if err != nil {
		return nil, err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if statusOK := response.StatusCode >= 200 && response.StatusCode < 300; !statusOK {
		return nil, fmt.Errorf("Non-OK HTTP status: %d (%s)", response.StatusCode, request.URL)
	}

	return response, nil
}