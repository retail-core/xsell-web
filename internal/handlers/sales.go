package handlers

import (
	"net/http"

	"github.com/retail-core/xsell-web/internal/clients"
)

type SalesHandler struct {
	ShellHandler
	SalesClient *clients.SalesClient
}


type SaleOrder struct {
	StaffName, Time, Amount, PaymentMethod string
}

type SaleGroup struct {
	Label      string
	OrderCount int
	Total      string
	Orders     []SaleOrder
}

func getSalesData() map[string]any {
	return map[string]any{
		"TodayProfit":   "₦12,090",
		"TodayRevenue":  "₦18,200",
		"ProfitTrend":   "+23%",
		"RevenueTrend":  "+15%",
		"TodayOrders":   2,
		"Groups": []SaleGroup{
			{
				Label: "Today", OrderCount: 2, Total: "₦60,300",
				Orders: []SaleOrder{
					{"Obimba Smart", "Aug 28th, 16:22:40", "₦44,200", "Cash"},
					{"Obimba Smart", "Aug 28th, 16:22:16", "₦16,100", "Transfer"},
				},
			},
			{
				Label: "Aug 2026", OrderCount: 6, Total: "₦87,350",
				Orders: []SaleOrder{
					{"Garry Lukeman", "Aug 2nd, 10:35:43", "₦35,600", "Transfer"},
					{"Obimba Smart", "Aug 1st, 23:29:14", "₦8,750", "Transfer"},
					{"Obimba Smart", "Aug 1st, 22:54:35", "₦23,550", "Transfer"},
				},
			},
		},
		"ContentTemplate": "sales-content",
	}
}

func (h *SalesHandler) Sales(w http.ResponseWriter, r *http.Request) {
	h.renderShell(w, "sales", "Sales",
		[]HeaderAction{{Icon: "filter", OnClick: "toggleFilter()"}, {Icon: "search", OnClick: "toggleSearch()"}},
		getSalesData(),
	)
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

func getOrderDetail(orderID string) OrderDetail {
	return OrderDetail{
		OrderID: "#233EFNI41",
		SoldBy:  "Obimba Nkechi",
		DateTime: "Feb 29th, 2026 · 13:00",
		Items: []OrderLineItem{
			{"Beef Crowich", "", "₦7,817", "₦23,451", 3},
			{"Chicken Wrap", "", "₦8,200", "₦8,200", 1},
		},
		Subtotal:      "₦31,651",
		PaymentMethod: "Cash",
		Status:        "Credit",
		IsCredit:      true,
		CreditToName:  "Rabi Aliyu",
	}
}