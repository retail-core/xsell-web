package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ProductRow struct {
	ID, Name, Category, Price, ImageURL string
	StockCount                          int
	StockLabel                          string // "24 items", "0 items"
	LowStock                            bool   // drives amber/red styling
}

func getProductsData() map[string]any {
	categories := []Category{
		{"All", "589 items", true},
		{"Spirit & Liquors", "201 items", false},
		{"Beer", "93 items", false},
		{"Wine", "322 items", false},
		{"Drinks", "376 items", false},
	}
	products := []ProductRow{
		{ID: "p-1", Name: "44cl Nuts Cashew Lion Baked 100g", Category: "Snacks", Price: "₦1,200", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10", StockCount: 33, StockLabel: "33 in stock", LowStock: false},
		{ID: "p-2", Name: "Energy Drink Regular 250ml", Category: "Drinks", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-3", Name: "Chocolate Mini Snickers", Category: "Personal Care", Price: "₦2,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s", StockCount: 5, StockLabel: "6 in stock", LowStock: true},
		{ID: "p-4", Name: "Noodles Indomie Chicken", Category: "Groceries", Price: "₦400", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-5", Name: "Chips Pringles 165g", Category: "Noodles", Price: "₦1,600", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10", StockCount: 33, StockLabel: "0 in stock", LowStock: true},
		{ID: "p-6", Name: "Noodles Chicken Honey", Category: "Milk", Price: "₦450", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-7", Name: "Cooking Oil Vegetable 1L", Category: "Detergent", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-8", Name: "44cl Nuts Cashew Lion Baked 100g", Category: "Snacks", Price: "₦1,200", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10", StockCount: 33, StockLabel: "33 in stock", LowStock: false},
		{ID: "p-9", Name: "Energy Drink Regular 250ml", Category: "Drinks", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-10", Name: "Chocolate Mini Snickers", Category: "Personal Care", Price: "₦2,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s", StockCount: 5, StockLabel: "6 in stock", LowStock: true},
		{ID: "p-11", Name: "Noodles Indomie Chicken", Category: "Groceries", Price: "₦400", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-12", Name: "Chips Pringles 165g", Category: "Noodles", Price: "₦1,600", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10", StockCount: 33, StockLabel: "0 in stock", LowStock: true},
		{ID: "p-13", Name: "Noodles Chicken Honey", Category: "Milk", Price: "₦450", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-14", Name: "Cooking Oil Vegetable 1L", Category: "Detergent", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-15", Name: "Croissant Maxi 15 Pack", Category: "Spices & Seasonings", Price: "₦500", ImageURL: "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-16", Name: "44cl Nuts Cashew Lion Baked 100g", Category: "Snacks", Price: "₦1,200", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPlgO3mgTCURQB5N0_NrrLEUypENVzdmAC0sOdawqHMqDFELVCjbYxhuBI&s=10", StockCount: 33, StockLabel: "33 in stock", LowStock: false},
		{ID: "p-17", Name: "Energy Drink Regular 250ml", Category: "Drinks", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS-0oT8tYOMTx0nqNOIUwpxTyoBuVfqr_QC0j2EjlhRMw&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-18", Name: "Chocolate Mini Snickers", Category: "Personal Care", Price: "₦2,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCJpLvI8W0zoZcmXOpztZ2oHV2icIO_7aLytsEaZrn3w&s", StockCount: 5, StockLabel: "6 in stock", LowStock: true},
		{ID: "p-19", Name: "Noodles Indomie Chicken", Category: "Groceries", Price: "₦400", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaNn_MHV0VelIG728ClSQiNinBllBEuetCY61Btb6WUg&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-20", Name: "Chips Pringles 165g", Category: "Noodles", Price: "₦1,600", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTFgCeJx0XUKF3TfkcTxPW9Vo3-aZywwR2Uhx6qoLBrJlHLupC_hRMvdOu-&s=10", StockCount: 33, StockLabel: "0 in stock", LowStock: true},
		{ID: "p-21", Name: "Noodles Chicken Honey", Category: "Milk", Price: "₦450", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZRKlJIUNKnVS0z5hk40ZQXnsVXdCpvIpLoVCPFDYHXygr1mxb6kpTMM4&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-22", Name: "Cooking Oil Vegetable 1L", Category: "Detergent", Price: "₦9,000", ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRluxYA9GMH2lDN4UUmy0zWsQWcS9YeyBXCCch6sdo58q4HagoeDIKVY_g&s=10", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-23", Name: "Croissant Maxi 15 Pack", Category: "Spices & Seasonings", Price: "₦500", ImageURL: "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
		{ID: "p-24", Name: "Croissant Maxi 15 Pack", Category: "Spices & Seasonings", Price: "₦500", ImageURL: "https://vegetalfoods.com/wp-content/uploads/2021/04/4P8A0760.jpg", StockCount: 33, StockLabel: "6 in stock", LowStock: false},
	}
	return map[string]any{
		"Categories":      categories,
		"Products":        products,
		"ContentTemplate": "products-content",
	}
}

func (h *InventoryHandler) Products(w http.ResponseWriter, r *http.Request) {
	h.renderShell(w, "products", "Products",
		[]HeaderAction{{Icon: "search", OnClick: "toggleSearch()"}},
		getProductsData(),
	)
}

type ProductForm struct {
	ID, Name, Category, CostPrice, SellingPrice, Barcode string
	Quantity, ReorderLevel                               int
	ImageURL                                             string
	IsActive                                             bool
}

func getProductForm(id string) ProductForm {
	// hardcoded for now — real fetch from inventory-service later
	if id == "" {
		return ProductForm{IsActive: true}
	}
	return ProductForm{
		ID: id, Name: "Hennessy XO Cognac 750ml", Category: "Spirit & Liquors",
		CostPrice: "56000", SellingPrice: "67000", Barcode: "1289612391292",
		Quantity: 24, ReorderLevel: 5, ImageURL: "", IsActive: true,
	}
}

func (h (*InventoryHandler))  ProductFormPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	form := getProductForm(id)
	isEdit := id != ""

	var title = "New Product"
	if isEdit {
		title = "Edit Product"
	}

	h.Tmpl.ExecuteTemplate(w, "product-form-page", map[string]any{
		"Form":       form,
		"IsEdit":     isEdit,
		"Categories": []string{"Spirit & Liquors", "Beer", "Wine", "Snacks", "Drinks", "Groceries"},
		"PageTitle": title,
	})
}
