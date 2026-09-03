package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SalesClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewSalesClient(baseURL string) *SalesClient {
	return &SalesClient{BaseURL: baseURL, HTTP: newHTTPClient()}
}

type Order struct {
	ID            string  `json:"id"`
	CustomerName  *string `json:"customer_name"`
	TotalAmount   float64 `json:"total_amount"`
	TotalCost     float64 `json:"total_cost"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	SoldBy        string  `json:"sold_by"`
	CreatedAt     string  `json:"created_at"`
}

func (c *SalesClient) GetOrders(storeID, token string) ([]Order, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/stores/"+storeID+"/orders", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sales service unreachable: %w", err)
	}
	defer resp.Body.Close()

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

type CreateOrderRequest struct {
	PaymentMethod  string                      `json:"payment_method"`
	InventoryItems []InventoryOrderItemRequest `json:"inventory_items"`
	SoldBy         *string                     `json:"sold_by,omitempty"`
}

type InventoryOrderItemRequest struct {
	InventoryID string `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
}

func (c *SalesClient) CreateOrder(storeID, token string, req CreateOrderRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/stores/"+storeID+"/orders", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("sales service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("order creation failed with status %d", resp.StatusCode)
	}
	return nil
}

type OrderItem struct {
	InventoryID string  `json:"inventory_id"`
	ProductName string  `json:"product_name"`
	ImageUrl    *string `json:"image_url"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderDetail struct {
	Order
	ReceiptNo *string     `json:"receipt_no,omitempty"`
	Items     []OrderItem `json:"items"`
}

func (c *SalesClient) GetOrderDetail(storeID, orderID, token string) (*OrderDetail, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/stores/"+storeID+"/orders/"+orderID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sales service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("order fetch failed with status %d", resp.StatusCode)
	}

	var detail OrderDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}
	return &detail, nil
}