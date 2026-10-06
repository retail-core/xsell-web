package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/session"
)

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

func (h *AccountHandler) StaffPage(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r)
	activeStore, _ := session.GetActiveStore(r)

	staff, err := h.AccountClient.GetStaff(activeStore.ID, token)
	if err != nil {
		h.Log.Error("failed to fetch staff", "store_id", activeStore.ID, "error", err)
		h.Tmpl.ExecuteTemplate(w, "staff-page", map[string]any{"LoadError": true})
		return
	}

	rows := buildStaffRows(staff)

	h.Tmpl.ExecuteTemplate(w, "staff-page", map[string]any{
		"Staff": rows,
		"PageTitle": "My Staff",
	})
}

type StaffRow struct {
	ID, Name, Role, Status string
}

func buildStaffRows(staff []clients.Staff) []StaffRow {
	rows := make([]StaffRow, 0, len(staff))
	for _, s := range staff {
		status := "Pending"
		if s.IsVerified {
			status = "Active"
		}
		rows = append(rows, StaffRow{
			ID:     s.ID,
			Name:   s.Username,
			Role:   s.Role,
			Status: status,
		})
	}
	return rows
}