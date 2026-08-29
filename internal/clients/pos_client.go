package clients

import (
	"net/http"
)

type PosClient struct {
	BaseURL string
	HTTP    *http.Client
}


func NewPosClient(baseURL string) *PosClient {
	return &PosClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}
