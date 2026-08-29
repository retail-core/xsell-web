package handlers

import "net/http"

type Product struct {
	Name, Price, ImageURL string
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

func getCartData() CartData {
	return CartData{
		OrderNumber: "#345348",
		Items: []CartItem{
			{"Nuts Cashew Lion 100g", "₦1,200", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10", 3, "₦12,000"},
			{"Energy Drink Regular 250ml", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10", 5, "₦12,000"},
			{"Chocolate Mini Snickers", "₦2,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s", 2, "₦5,000"},
			{"Noodles Indomie Chicken", "₦400", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10", 4, "₦12,000"},
			{"Uht Soya Milk Strawberry Vitamilk 300ml", "₦1,000", "", 12, "₦12,000"},
			{"Nuts Cashew Lion 100g", "₦1,200", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10", 3, "₦12,000"},
			{"Energy Drink Regular 250ml", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10", 5, "₦12,000"},
			{"Chocolate Mini Snickers", "₦2,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s", 2, "₦5,000"},
			{"Noodles Indomie Chicken", "₦400", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10", 4, "₦12,000"},
			{"Uht Soya Milk Strawberry Vitamilk 300ml", "₦1,000", "", 12, "₦12,000"},
		},
		Total: "₦25,800",
	}
}

var categories = []Category{
	{"All", "1,201 items", true},
	{"Beer", "93 items", false},
	{"Wine", "322 items", false},
	{"Snacks", "410 items", false},
	{"Drinks", "376 items", false},
}

var products = []Product{
	{"Nuts Cashew Lion 100g", "₦1,200", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10"},
	{"Energy Drink Regular 250ml", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10"},
	{"Chocolate Mini Snickers", "₦2,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s"},
	{"Noodles Indomie Chicken", "₦400", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10"},
	{"Chips Pringles 165g", "₦1,600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10"},
	{"Noodles Chicken Honey", "₦450", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10"},
	{"Cooking Oil Vegetable 1L", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10"},
	{"Croissant Maxi 15 Pack", "₦500", "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg"},
	{"French Loaf Large 200g", "₦600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRdNaxyms_nlxh_xM-cXwL53ZH4NLKu6eGkHKLIJeFKjfXe9yBL2oFjSOdu&s=10"},
	{"Nuts Cashew Lion 100g", "₦1,200", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10"},
	{"Energy Drink Regular 250ml", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10"},
	{"Chocolate Mini Snickers", "₦2,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s"},
	{"Noodles Indomie Chicken", "₦400", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10"},
	{"Chips Pringles 165g", "₦1,600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10"},
	{"Noodles Chicken Honey", "₦450", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10"},
	{"Cooking Oil Vegetable 1L", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10"},
	{"Croissant Maxi 15 Pack", "₦500", "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg"},
	{"French Loaf Large 200g", "₦600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRdNaxyms_nlxh_xM-cXwL53ZH4NLKu6eGkHKLIJeFKjfXe9yBL2oFjSOdu&s=10"},
		{"Nuts Cashew Lion 100g", "₦1,200", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10"},
	{"Energy Drink Regular 250ml", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10"},
	{"Chocolate Mini Snickers", "₦2,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s"},
	{"Noodles Indomie Chicken", "₦400", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10"},
	{"Chips Pringles 165g", "₦1,600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10"},
	{"Noodles Chicken Honey", "₦450", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10"},
	{"Cooking Oil Vegetable 1L", "₦9,000", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10"},
	{"Croissant Maxi 15 Pack", "₦500", "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg"},
	{"French Loaf Large 200g", "₦600", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRdNaxyms_nlxh_xM-cXwL53ZH4NLKu6eGkHKLIJeFKjfXe9yBL2oFjSOdu&s=10"},
}

func (h *InventoryHandler) POS(w http.ResponseWriter, r *http.Request) {
	h.renderShell(w, "pos", "Point of Sale", []HeaderAction{
		{Icon: "scan", OnClick: ""},
		{Icon: "search", OnClick: ""},
	},
		map[string]any{
			"Categories": categories,
			"Products":   products,
			"CartCount":  3,
			"ContentTemplate": "pos-content",
			"Cart": getCartData(),
		},
	)
}
