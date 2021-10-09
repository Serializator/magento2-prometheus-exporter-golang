package collector

type environmentInfoResponse struct {
	Version string `json:"version"`
	Edition string `json:"edition"`
	Mode    string `json:"mode"`
}
