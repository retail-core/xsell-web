package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type ProductClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}

func (c *ProductClient) UploadImage(fileHeader *multipart.FileHeader, token string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image_file", fileHeader.Filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+"/products/images/upload", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("product service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image upload failed with status %d", resp.StatusCode)
	}

	var result struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ImageURL, nil
}