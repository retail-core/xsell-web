package handlers

import (
	"net/http"

	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/session"
	m_middleware "github.com/retail-core/xsell-web/internal/middleware"
)

type Inventory struct {
	ID, Name, Price, ImageURL string
	SellingpPrice float32
}

type Category struct {
	Name, Count string
	Active      bool
}

type CartItem struct {
	Name, Price, ImageURL string
	Qty                   int
	LineTotal              string
}

type CartData struct {
	OrderNumber string
	Items       []CartItem
	Total       string
}

func (h *InventoryHandler) POS(w http.ResponseWriter, r *http.Request) {
	token := m_middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)

	if activeStore == nil {
		w.Header().Set("X-Toast-Message", "No active store selected")
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	items, err := h.InventoryClient.GetInventory(activeStore.ID, token)
	if err != nil {
		h.Log.Error("failed to fetch inventory for POS", "store_id", activeStore.ID, "error", err)
		h.renderShell(w, r, "pos", "Point of Sale", nil, map[string]any{
			"ContentTemplate": "pos-content",
			"LoadError":       true,
		})
		return
	}

	// POS only shows sellable (active, in-stock) items
	activeItems := make([]clients.InventoryItem, 0, len(items))
	for _, item := range items {
		if item.IsActive {
			activeItems = append(activeItems, item)
		}
	}

	categories := buildCategoryChips(activeItems)
	products := buildPOSProducts(activeItems)

	h.renderShell(w, r, "pos", "Point of Sale",
		[]HeaderAction{{Icon: "search", OnClick: "openPageSearch()"}, {Icon: "scan", OnClick: "#"}},
		map[string]any{
			"ContentTemplate": "pos-content",
			"Categories":      categories,
			"Products":        products,
			"CartCount":       0,
		},
	)
}

type Product struct {
	ID, Name, ImageURL, Category string
	Price               string  // formatted, for display
	PriceRaw            float64 // raw, for cart math
}

func buildPOSProducts(items []clients.InventoryItem) []Product {
	products := make([]Product, 0, len(items))
	for _, item := range items {
		if item.TotalQty  == 0 {
			continue
		}
		
		name := "Unnamed Product"
		if item.Name != nil {
			name = *item.Name
		}
		imageURL := ""
		if item.ImageUrl != nil {
			imageURL = *item.ImageUrl
		}
		products = append(products, Product{
			ID:       item.ID,
			Name:     name,
			Price:    formatNaira(item.SellingPrice),
			PriceRaw: item.SellingPrice,
			ImageURL: imageURL,
			Category: *item.Category,
		})
	}
	return products
}