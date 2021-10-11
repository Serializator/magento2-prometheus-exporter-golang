package collector

type ordersResponse struct {
	Items []struct {
		StoreId int64 `json:"store_id"`
		State  string `json:"state"`
		Status string `json:"status"`
		Payment struct{
			Method string `json:"method"`
		} `json:"payment"`
	} `json:"items"`
}
