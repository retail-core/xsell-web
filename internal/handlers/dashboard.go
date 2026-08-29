package handlers

import "net/http"

// type NavItem struct {
// 	Href, Icon, Label string
// 	Active             bool
// }

func (h *AccountHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	navItems := []NavItem{
		{"/pos", "📊", "Overview", true},
		{"/sales", "₦", "Sales", false},
		{"/products", "📦", "Products", false},
		{"/customers", "👥", "Customers", false},
		{"/reports", "📈", "Reports", false},
		{"/settings", "⚙️", "Settings", false},
	}
	h.Tmpl.ExecuteTemplate(w, "base", map[string]any{
		"NavItems": navItems,
		"Active":   "pos",
	})
}