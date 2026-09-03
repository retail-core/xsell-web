package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"
	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/session"
)

type ProductForm struct {
	ID, Name, Category, CostPrice, SellingPrice, Barcode string
	Quantity, ReorderLevel                               int
	ImageURL                                             string
	IsActive                                             bool
}

func (h *InventoryHandler) ProductFormPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	isEdit := id != ""
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)

	var form ProductForm
	if isEdit {
		items, err := h.InventoryClient.GetInventory(activeStore.ID, token) // same call as list page
		if err != nil {
			h.Log.Error("failed to fetch inventory", "error", err)
			w.Header().Set("X-Toast-Message", "Could not load product")
			http.Redirect(w, r, "/products", http.StatusSeeOther)
			return
		}

		var found *clients.InventoryItem
		for _, item := range items {
			if item.ID == id {
				found = &item
				break
			}
		}
		if found == nil {
			w.Header().Set("X-Toast-Message", "Product not found")
			http.Redirect(w, r, "/products", http.StatusSeeOther)
			return
		}
		form = mapInventoryToForm(found)
	}

	var title = "New Product"
	if isEdit {
		title = "Edit Product"
	}

	err := h.Tmpl.ExecuteTemplate(w, "product-form-page", map[string]any{
		"Form":       form,
		"IsEdit":     isEdit,
		"PageTitle":  title,
	})

	if err != nil {
		h.Log.Error("error loading page", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) Products(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	role := middleware.GetRole(r)

	if activeStore == nil {
		w.Header().Set("X-Toast-Message", "No active store selected")
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	items, err := h.InventoryClient.GetInventory(activeStore.ID, token)
	if err != nil {
		h.Log.Error("failed to fetch inventory", "store_id", activeStore.ID, "error", err)
		h.renderShell(w, r, "products", "Inventory", nil, map[string]any{
			"ContentTemplate": "products-content",
			"LoadError":       true,
		})
		return
	}

	categories := buildCategoryChips(items)
	rows := buildProductRows(items)

	h.renderShell(w, r, "products", "Inventory",
		[]HeaderAction{{Icon: "search", OnClick: "openPageSearch()"}},
		map[string]any{
			"ContentTemplate": "products-content",
			"Categories":      categories,
			"Products":        rows,
			"IsOwner":         role == "business_owner",
		},
	)
}

type CategoryChip struct {
	Name   string
	Count  int
	Active bool
}

func buildCategoryChips(items []clients.InventoryItem) []CategoryChip {
	counts := map[string]int{}
	for _, item := range items {
		cat := "Uncategorized"
		if item.Category != nil && *item.Category != "" {
			cat = *item.Category
		}
		counts[cat]++
	}

	chips := []CategoryChip{{Name: "All", Count: len(items), Active: true}}

	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		chips = append(chips, CategoryChip{Name: name, Count: counts[name]})
	}
	return chips
}

type ProductRow struct {
	ID, Name, Category, Price, ImageURL string
	StockCount                          int
	StockLabel                          string
	LowStock                            bool
	SellingPriceRaw, CostPriceRaw       float64
}

func buildProductRows(items []clients.InventoryItem) []ProductRow {
	rows := make([]ProductRow, 0, len(items))
	for _, item := range items {
		name := "Unnamed Product"
		if item.Name != nil {
			name = *item.Name
		}
		category := "Uncategorized"
		if item.Category != nil {
			category = *item.Category
		}
		imageURL := ""
		if item.ImageUrl != nil {
			imageURL = *item.ImageUrl
		}

		lowStock := false
		if item.MinThreshold != nil && item.TotalQty <= *item.MinThreshold {
			lowStock = true
		}

		costPrice := 0.0
		if item.CostPrice != nil {
			costPrice = *item.CostPrice
		}

		rows = append(rows, ProductRow{
			ID:              item.ID,
			Name:            name,
			Category:        category,
			Price:           formatNaira(item.SellingPrice),
			ImageURL:        imageURL,
			StockCount:      item.TotalQty,
			StockLabel:      fmt.Sprintf("%d in stock", item.TotalQty),
			LowStock:        lowStock,
			SellingPriceRaw: item.SellingPrice,
			CostPriceRaw:    costPrice,
		})
	}
	return rows
}

func formatNaira(amount float64) string {
	return fmt.Sprintf("₦%s", humanizeNumber(amount))
}

func (h *InventoryHandler) RestockSubmit(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	inventoryID := chi.URLParam(r, "id")

	var body struct {
		Quantity     int     `json:"quantity"`
		SellingPrice float64 `json:"selling_price"`
		CostPrice    float64 `json:"cost_price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	req := clients.RestockRequest{
		Quantity:     body.Quantity,
		SellingPrice: body.SellingPrice,
		CostPrice:    &body.CostPrice,
		Type:         "set",
	}

	if err := h.InventoryClient.Restock(activeStore.ID, inventoryID, token, req); err != nil {
		h.Log.Error("restock failed", "inventory_id", inventoryID, "error", err)
		w.Header().Set("X-Toast-Message", "Could not update stock")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Toast-Message", "Stock updated")
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	inventoryID := chi.URLParam(r, "id")

	if err := h.InventoryClient.DeleteInventory(activeStore.ID, inventoryID, token); err != nil {
		h.Log.Error("delete inventory failed", "inventory_id", inventoryID, "error", err)
		w.Header().Set("X-Toast-Message", "Could not delete product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Toast-Message", "Product deleted")
	w.WriteHeader(http.StatusOK)
}

func mapInventoryToForm(item *clients.InventoryItem) ProductForm {
	name, category, barcode, imageURL := "", "", "", ""
	if item.Name != nil {
		name = *item.Name
	}
	if item.Category != nil {
		category = *item.Category
	}
	if item.Barcode != nil {
		barcode = *item.Barcode
	}
	if item.ImageUrl != nil {
		imageURL = *item.ImageUrl
	}

	costPrice := 0.0
	if item.CostPrice != nil {
		costPrice = *item.CostPrice
	}

	return ProductForm{
		ID:           item.ID,
		Name:         name,
		Category:     category,
		Barcode:      barcode,
		ImageURL:     imageURL,
		CostPrice:    fmt.Sprintf("%.2f", costPrice),
		SellingPrice: fmt.Sprintf("%.2f", item.SellingPrice),
		Quantity:     item.TotalQty,
		IsActive:     item.IsActive,
	}
}
