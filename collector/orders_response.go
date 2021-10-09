package collector

type ordersResponse struct {
	Items []struct {
		State  string `json:"state"`
		Status string `json:"status"`
	} `json:"items"`
}
