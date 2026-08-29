package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store struct {
	ID, Name, Address string
	IsActive           bool
}

func getStoresData() []Store {
	return []Store{
		{"1", "Kilimanjaro Adubawa Enyean", "Benin, Edo", true},
		{"2", "Rite Food Restaurant", "Port Harcourt, Rivers", true},
		{"3", "Vivid Gaming Experience", "Lagos, Lagos", true},
		{"4", "Chase Haven Pharmacy", "Abuja, FCT", false},
	}
}

func (h *AccountHandler) StoresPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "stores-page", map[string]any{
		"Stores": getStoresData(),
	})
}

type StoreForm struct {
	ID, Name, Address string
	IsActive           bool
}

func getStoreForm(id string) StoreForm {
	if id == "" {
		return StoreForm{IsActive: true}
	}
	return StoreForm{
		ID: id, Name: "Kilimanjaro Adubawa Enyean", Address: "Benin, Edo", IsActive: true,
	}
}

func (h *AccountHandler) StoreFormPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	form := getStoreForm(id)
	isEdit := id != ""

	title := "New Store"
	submitLabel := "Create Store"
	if isEdit {
		title = "Edit Store"
		submitLabel = "Save Changes"
	}

	h.Tmpl.ExecuteTemplate(w, "store-form-page", map[string]any{
		"Form":        form,
		"IsEdit":      isEdit,
		"PageTitle":   title,
		"SubmitLabel": submitLabel,
	})
}