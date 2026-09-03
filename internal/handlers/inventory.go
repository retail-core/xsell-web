package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/session"
)

type InventoryHandler struct {
	ShellHandler
	InventoryClient *clients.InventoryClient
}

func (h *InventoryHandler) ProductFormSubmit(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	id := chi.URLParam(r, "id")
	isEdit := id != ""

	r.ParseForm()
	name := r.FormValue("name")
	category := r.FormValue("category")
	barcode := r.FormValue("barcode")

	costPrice, _ := strconv.ParseFloat(r.FormValue("cost_price"), 64)
	sellingPrice, _ := strconv.ParseFloat(r.FormValue("selling_price"), 64)
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))
	reorderLevel, _ := strconv.Atoi(r.FormValue("reorder_level"))

	if name == "" || category == "" {
		w.Header().Set("X-Toast-Message", "Name and category are required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := clients.CreateInventoryRequest{
		Name:         &name,
		Category:     &category,
		SellingPrice: sellingPrice,
		CostPrice:    &costPrice,
		MinThreshold: &reorderLevel,
		InitialQty:   quantity,
	}
	if barcode != "" {
		req.Barcode = &barcode
	}

	var err error
	if isEdit {
		// err = h.InventoryClient.UpdateInventory(activeStore.ID, id, token, req)
	} else {
		err = h.InventoryClient.CreateInventory(activeStore.ID, token, req)
	}

	if err != nil {
		h.Log.Error("save product failed", "is_edit", isEdit, "error", err)
		w.Header().Set("X-Toast-Message", "Could not save product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Toast-Message", "Product saved")
	w.WriteHeader(http.StatusOK)
}