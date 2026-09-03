package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type NavItem struct {
	Href, Icon, Label string
	Active            bool
}

type HeaderAction struct {
	Icon    string
	OnClick template.JS
}

// Shared by every handler that renders inside the app shell (POS, Products, Sales, Account, etc.)
type ShellHandler struct {
	Tmpl *template.Template
	Log  *slog.Logger
}

func (s *ShellHandler) renderShell(w http.ResponseWriter, r *http.Request, active, title string, actions []HeaderAction, extra map[string]any) {
	navItems := []NavItem{
		{"/pos", "cart", "POS", active == "pos"},
		{"/sales", "cash", "Sales", active == "sales"},
		{"/products", "cube", "Products", active == "products"},
		{"/account", "user", "More", active == "account"},
	}

	data := map[string]any{
		"NavItems":      navItems,
		"Active":        active,
		"PageTitle":     title,
		"HeaderActions": actions,
	}
	for k, v := range extra {
		data[k] = v
	}

	if r.Header.Get("HX-Request") == "true" {
		// Main content
		s.Tmpl.ExecuteTemplate(w, data["ContentTemplate"].(string), data)

		// Out-of-band swap: topbar
		fmt.Fprint(w, `<div id="mobile-topbar-wrapper" hx-swap-oob="true">`)
		s.Tmpl.ExecuteTemplate(w, "mobile-topbar", data)
		fmt.Fprint(w, `</div>`)

		// Out-of-band swap: bottom nav
		fmt.Fprint(w, `<div id="bottom-nav-wrapper" hx-swap-oob="true">`)
		s.Tmpl.ExecuteTemplate(w, "bottom-nav", data)
		fmt.Fprint(w, `</div>`)

		return
	}

	err := s.Tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		s.Log.Error("error loading page", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
