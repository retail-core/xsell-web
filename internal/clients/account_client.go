package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountClient struct {
	BaseURL string
	HTTP    *http.Client
}


func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}

type Store struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address *string `json:"address,omitempty"`
	Tag     *string `json:"tag,omitempty"`
}

func (c *AccountClient) GetUserStores(userID, token string) ([]Store, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/v1/users/"+userID+"/stores", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("account service unreachable: %w", err)
	}
	defer resp.Body.Close()

	var stores []Store
	if err := json.NewDecoder(resp.Body).Decode(&stores); err != nil {
		return nil, err
	}
	return stores, nil
}