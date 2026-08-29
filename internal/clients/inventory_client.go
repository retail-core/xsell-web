package clients

import (
	"net/http"
)

type InventoryClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}