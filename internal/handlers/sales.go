package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/retail-core/xsell-web/internal/clients"
	m_middleware "github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/session"
)

type SalesHandler struct {
	ShellHandler
	SalesClient *clients.SalesClient
}

type OrderLineItem struct {
	Name, ImageURL, UnitPrice, LineTotal string
	Qty                                  int
}

type OrderDetail struct {
	OrderID, SoldBy, DateTime            string
	Items                                []OrderLineItem
	Subtotal, PaymentMethod, Status      string
	IsCredit                             bool
	CreditToName                         string
}

func (h *SalesHandler) Sales(w http.ResponseWriter, r *http.Request) {
	token := m_middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)

	if activeStore == nil {
		w.Header().Set("X-Toast-Message", "No active store selected")
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	orders, err := h.SalesClient.GetOrders(activeStore.ID, token)
	if err != nil {
		h.Log.Error("failed to fetch orders", "store_id", activeStore.ID, "error", err)
		h.renderShell(w, r, "sales", "Sales", nil, map[string]any{
			"ContentTemplate": "sales-content",
			"LoadError":       true,
		})
		return
	}

	todayProfit, todayRevenue, todayCount := computeTodaySummary(orders)
	groups := groupOrdersByPeriod(orders)

	h.renderShell(w, r, "sales", "Sales",
		[]HeaderAction{{Icon: "filter", OnClick: "toggleFilter()"}, {Icon: "search", OnClick: "toggleSearch()"}},
		map[string]any{
			"ContentTemplate": "sales-content",
			"TodayProfit":     formatNaira(todayProfit),
			"TodayRevenue":    formatNaira(todayRevenue),
			"TodayOrders":     todayCount,
			"Groups":          groups,
		},
	)
}

type SaleOrder struct {
	ID, StaffName, Time, Amount, PaymentMethod string
	Status                                     string
}

type SaleGroup struct {
	Label      string
	OrderCount int
	Total      string
	Orders     []SaleOrder
}

func computeTodaySummary(orders []clients.Order) (profit, revenue float64, count int) {
	now := time.Now()
	for _, o := range orders {
		t, err := time.Parse(time.RFC3339, o.CreatedAt)
		if err != nil {
			continue
		}
		if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
			revenue += o.TotalAmount
			profit += o.TotalAmount - o.TotalCost
			count++
		}
	}
	return
}

func groupOrdersByPeriod(orders []clients.Order) []SaleGroup {
	type groupKey struct {
		label string
		sortKey string // for ordering groups newest-first
	}

	grouped := map[string][]clients.Order{}
	labelOrder := []string{}
	seen := map[string]bool{}

	now := time.Now()

	for _, o := range orders {
		t, err := time.Parse(time.RFC3339, o.CreatedAt)
		if err != nil {
			continue
		}

		var label string
		if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
			label = "Today"
		} else if t.Year() == now.Year() && t.YearDay() == now.YearDay()-1 {
			label = "Yesterday"
		} else {
			label = t.Format("January 2006") // e.g. "Aug 2026"
		}

		grouped[label] = append(grouped[label], o)
		if !seen[label] {
			seen[label] = true
			labelOrder = append(labelOrder, label)
		}
	}

	// Sort orders within each group newest-first
	for label := range grouped {
		sort.Slice(grouped[label], func(i, j int) bool {
			ti, _ := time.Parse(time.RFC3339, grouped[label][i].CreatedAt)
			tj, _ := time.Parse(time.RFC3339, grouped[label][j].CreatedAt)
			return ti.After(tj)
		})
	}

	groups := make([]SaleGroup, 0, len(labelOrder))
	for _, label := range labelOrder {
		orders := grouped[label]
		var total float64
		rows := make([]SaleOrder, 0, len(orders))
		for _, o := range orders {
			total += o.TotalAmount
			t, _ := time.Parse(time.RFC3339, o.CreatedAt)
			customerLabel := "Walk-in Customer"
			if o.CustomerName != nil && *o.CustomerName != "" {
				customerLabel = *o.CustomerName
			}
			rows = append(rows, SaleOrder{
				ID:            o.ID,
				StaffName:     o.SoldBy,
				Time:          t.Format("Jan 2, 15:04"),
				Amount:        formatNaira(o.TotalAmount),
				PaymentMethod: titleCase(o.PaymentMethod),
				Status:        o.Status,
			})
			_ = customerLabel // available if you want to show it per-row later
		}
		groups = append(groups, SaleGroup{
			Label:      label,
			OrderCount: len(orders),
			Total:      formatNaira(total),
			Orders:     rows,
		})
	}

	return groups
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]) + toLower(s[1:])
}
func toLower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}


func (h *SalesHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	token := m_middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	soldBy := session.GetUserName(r)


	if activeStore == nil {
		http.Error(w, "no active store", http.StatusBadRequest)
		return
	}

	var body clients.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}


	if soldBy != "" {
		body.SoldBy = &soldBy
	}


	if len(body.InventoryItems) == 0 {
		http.Error(w, "cart is empty", http.StatusBadRequest)
		return
	}

	err := h.SalesClient.CreateOrder(activeStore.ID, token, body)
	if err != nil {
		h.Log.Error("checkout failed", "store_id", activeStore.ID, "error", err)
		http.Error(w, "checkout failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesHandler) OrderDetail(w http.ResponseWriter, r *http.Request) {
	token := m_middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	orderID := chi.URLParam(r, "id")

	detail, err := h.SalesClient.GetOrderDetail(activeStore.ID, orderID, token)
	if err != nil {
		h.Log.Error("failed to fetch order detail", "order_id", orderID, "error", err)
		http.Error(w, "could not load order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}