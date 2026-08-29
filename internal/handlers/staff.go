package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StaffMember struct {
	ID, Name, Role, InitialsColor string
	Status                        string // "Active", "Pending", "Disabled"
}

func getStaffData() []StaffMember {
	return []StaffMember{
		{"1", "Garry Lukeman", "Sales Attendant", "accent", "Active"},
		{"2", "Amarachi Kelvin", "Store Manager", "accent", "Active"},
		{"3", "Rabi Aliyu", "Sales Attendant", "accent", "Pending"},
		{"4", "Chidi Okonkwo", "Cashier", "accent", "Disabled"},
	}
}

func (h *AccountHandler) StaffPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "staff-page", map[string]any{
		"Staff":     getStaffData(),
		"PageTitle": "My Staff",
	})
}

type StaffForm struct {
	ID, Name, Email, Role, Store string
}

func getStaffForm(id string) StaffForm {
	if id == "" {
		return StaffForm{}
	}
	return StaffForm{
		ID: id, Name: "Garry Lukeman", Email: "garry@example.com",
		Role: "Sales Attendant", Store: "Main Store",
	}
}

func (h *AccountHandler) StaffFormPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	form := getStaffForm(id)
	isEdit := id != ""

	title := "Add Staff"
	if isEdit {
		title = "Edit Staff"
	}

	submitLabel := "Send Invite"
	if isEdit {
		submitLabel = "Save Changes"
	}

	h.Tmpl.ExecuteTemplate(w, "staff-form-page", map[string]any{
		"Form":      form,
		"IsEdit":    isEdit,
		"PageTitle": title,
		"submitLabel": submitLabel,
		"Stores":    []string{"Main Store", "Rite Food Restaurant", "Vivid Gaming Experience", "Chase Haven Pharmacy"},
	})
}
