package clients

import (
	"net/http"
)

type SalesClient struct {
	BaseURL string
	HTTP    *http.Client
}


func NewSalesClient(baseURL string) *SalesClient {
	return &SalesClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}