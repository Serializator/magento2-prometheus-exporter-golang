package collector

type invoicesResponse struct {
	Total []struct {
		StoreId int    `json:"store_id"`
		State   string `json:"state"`
		Count   int    `json:"count"`
	} `json:"total"`
}
