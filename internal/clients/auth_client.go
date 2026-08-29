package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthClient struct {
	BaseURL string
	HTTP    *http.Client
}


func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         LUser  `json:"user"`
}

type LUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (c *AuthClient) Login(email, password string) (*LoginResponse, error) {
	payload, _ := json.Marshal(LoginRequest{Email: email, Password: password})

	req, err := http.NewRequest("POST", c.BaseURL+"/login", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth service unreachable: %w", err)
	}
	defer resp.Body.Close()

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	return &result, nil
}
