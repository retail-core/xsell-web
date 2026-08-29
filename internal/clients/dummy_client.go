package clients

import (
	"net/http"
)

type DummyClient struct {
	BaseURL string
	HTTP    *http.Client
}


func NewDummyClient(baseURL string) *DummyClient {
	return &DummyClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}