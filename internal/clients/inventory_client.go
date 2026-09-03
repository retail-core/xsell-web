package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InventoryClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}

type InventoryItem struct {
	ID           string   `json:"id"`
	StoreID      string   `json:"storeId"`
	ProductID    *string  `json:"productId"`
	Name         *string  `json:"name"`
	Category     *string  `json:"category"`
	Barcode      *string  `json:"barcode"`
	ImageUrl     *string  `json:"imageUrl"`
	SellingPrice float64  `json:"sellingPrice"`
	CostPrice    *float64 `json:"costPrice"`
	TotalQty     int      `json:"totalQty"`
	IsActive     bool     `json:"isActive"`
	MinThreshold *int     `json:"minThreshold"`
}

func (c *InventoryClient) GetInventory(storeID, token string) ([]InventoryItem, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/stores/"+storeID+"/inventories", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("inventory service unreachable: %w", err)
	}
	defer resp.Body.Close()

	var items []InventoryItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

type RestockRequest struct {
	Quantity     int      `json:"quantity"`
	SellingPrice float64  `json:"selling_price"`
	CostPrice    *float64 `json:"cost_price"`
	Type         string   `json:"type"`
}

func (c *InventoryClient) Restock(storeID, inventoryID, token string, req RestockRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/stores/"+storeID+"/inventories/"+inventoryID+"/restock", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("inventory service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("restock failed with status %d", resp.StatusCode)
	}
	return nil
}

func (c *InventoryClient) DeleteInventory(storeID, inventoryID, token string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+"/stores/"+storeID+"/inventories/"+inventoryID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("inventory service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete failed with status %d", resp.StatusCode)
	}
	return nil
}

type CreateInventoryRequest struct {
	Name         *string  `json:"name,omitempty"`
	Category     *string  `json:"category,omitempty"`
	Barcode      *string  `json:"barcode,omitempty"`
	SellingPrice float64  `json:"selling_price"`
	CostPrice    *float64 `json:"cost_price,omitempty"`
	MinThreshold *int     `json:"min_threshold,omitempty"`
	InitialQty   int      `json:"initial_qty"`
}

func (c *InventoryClient) CreateInventory(storeID, token string, req CreateInventoryRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/stores/"+storeID+"/inventories", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("inventory service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create failed with status %d", resp.StatusCode)
	}
	return nil
}

type UpdateInventoryStatusRequest struct {
	Status bool `json:"status"`
}

func (c *InventoryClient) UpdateStatus(storeID, inventoryID, token string, active bool) error {
	payload, _ := json.Marshal(UpdateInventoryStatusRequest{Status: active})

	req, err := http.NewRequest("PATCH", c.BaseURL+"/stores/"+storeID+"/inventories/"+inventoryID+"/status", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("inventory service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status update failed with status %d", resp.StatusCode)
	}
	return nil
}