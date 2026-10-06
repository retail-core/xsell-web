package handlers

import (
	"encoding/json"
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
	ProductClient   *clients.ProductClient
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
	isActive := r.FormValue("is_active") == "true"
	imageURL := r.FormValue("image_url")

	if name == "" || category == "" {
		w.Header().Set("X-Toast-Message", "Name and category are required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	if isEdit {
		req := clients.UpdateInventoryRequest{
			Name:         &name,
			Category:     &category,
			SellingPrice: &sellingPrice,
			CostPrice:    &costPrice,
			MinThreshold: &reorderLevel,
			InitialQty:   &quantity,
			IsActive:     &isActive,
			ImageUrl:     &imageURL,
		}
		if barcode != "" {
			req.Barcode = &barcode
		}

		if imageURL != "" {
			req.ImageUrl = &imageURL 
		}

		err = h.InventoryClient.UpdateInventory(activeStore.ID, id, token, req)
	} else {
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

		if imageURL != "" {
			req.ImageUrl = &imageURL 
		}

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

func (h *InventoryHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)
	inventoryID := chi.URLParam(r, "id")

	var body struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.InventoryClient.UpdateStatus(activeStore.ID, inventoryID, token, body.Active); err != nil {
		h.Log.Error("status update failed", "inventory_id", inventoryID, "error", err)
		w.Header().Set("X-Toast-Message", "Could not update product status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := "Product deactivated"
	if body.Active {
		msg = "Product activated"
	}
	w.Header().Set("X-Toast-Message", msg)
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)

	r.ParseMultipartForm(10 << 20) // 10MB max
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "no image provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := h.ProductClient.UploadImage(header, token)
	if err != nil {
		h.Log.Error("image upload failed", "error", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"image_url": ""})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"image_url": imageURL})
}
