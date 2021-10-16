package collector

type ordersResponse struct {
	Total []struct {
		StoreId int `json:"store_id"`
		Status string `json:"status"`
		State  string `json:"state"`
		PaymentMethod string `json:"payment_method"`
		Count int `json:"count"`
	} `json:"total"`
}
