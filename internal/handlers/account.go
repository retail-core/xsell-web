package handlers

import (
	"log"
	"net/http"

	"github.com/retail-core/xsell-web/internal/clients"
	m_middleware "github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/session"
)

type AccountHandler struct {
	ShellHandler
	AccountClient *clients.AccountClient
}

func (h *AccountHandler) Account(w http.ResponseWriter, r *http.Request) {
	activeStore, _ := session.GetActiveStore(r)
	role := m_middleware.GetRole(r)

	h.renderShell(w, r, "account", "Account", []HeaderAction{
		{Icon: "theme", OnClick: "toggleTheme()"},
		{Icon: "bell", OnClick: ""},
	}, map[string]any{
		"ActiveStore": activeStore,
		"ContentTemplate": "account-content",
		"IsOwner": role == "business_owner",
	})
}

func (h *AccountHandler) SubscriptionPage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "subscription-page", map[string]any{
		"Features": []string{
			"Unlimited sales & inventory tracking",
			"Multi-store management",
			"Staff accounts & access",
			"Customer credit & debt tracking",
			"Sales reports & business insights",
			"Barcode scanning",
		},
	})

	if err != nil {
		log.Println("template render error:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *AccountHandler) StoreSwitcherSheet(w http.ResponseWriter, r *http.Request) {
	userID := m_middleware.GetUserID(r)
	role := m_middleware.GetRole(r)
	token := m_middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)

	stores, err := h.AccountClient.GetUserStores(userID, token)
	if err != nil {
		h.Log.Error("failed to fetch stores for switcher", "user_id", userID, "error", err)
		w.Header().Set("X-Toast-Message", "Could not load stores")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.Info("stores fetched", "count", len(stores))

	var activeID string
	if activeStore != nil {
		activeID = activeStore.ID
	}

	h.Tmpl.ExecuteTemplate(w, "store-switcher-rows", map[string]any{
		"Stores":        stores,
		"ActiveStoreID": activeID,
		"IsOwner": role == "business_owner",
	})
}

func (h *AccountHandler) SwitchStore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	storeID := r.FormValue("store_id")
	userID := m_middleware.GetUserID(r)
	token := m_middleware.GetToken(r)

	stores, err := h.AccountClient.GetUserStores(userID, token)
	if err != nil {
		h.Log.Error("failed to fetch stores for switch", "user_id", userID, "error", err)
		w.Header().Set("X-Toast-Message", "Could not switch store")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.Info("stores fetched", "count", len(stores))

	var selected *clients.Store
	for _, s := range stores {
		if s.ID == storeID {
			selected = &s
			break
		}
	}
	if selected == nil {
		w.Header().Set("X-Toast-Message", "Store not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.SetActiveStore(w, *selected)
	w.Header().Set("X-Toast-Message", "Switched to "+selected.Name)
	w.Header().Set("HX-Refresh", "true") // full page reload to reflect new active store everywhere
	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) AccessDeniedPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "access-denied-page", nil)
}