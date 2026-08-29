package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type NavItem struct {
	Href, Icon, Label string
	Active            bool
}

type HeaderAction struct {
	Icon string
    OnClick template.JS
}

// Shared by every handler that renders inside the app shell (POS, Products, Sales, Account, etc.)
type ShellHandler struct {
	Tmpl *template.Template
	Log *slog.Logger
}

func (s *ShellHandler) renderShell(w http.ResponseWriter, active, title string, actions []HeaderAction, extra map[string]any) {
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

	err := s.Tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}