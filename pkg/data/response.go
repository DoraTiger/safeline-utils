package data



type CertListResponse struct {
	Data struct {
		Nodes []interface{} `json:"nodes"`
	} `json:"data"`
}

type UpdateCertResponse struct {
	Err *string `json:"err"`
}