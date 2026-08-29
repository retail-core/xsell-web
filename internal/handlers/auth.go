package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/session"
)

type AuthHandler struct {
	Tmpl       *template.Template
	AuthClient *clients.AuthClient
	AccountClient *clients.AccountClient
	Log *slog.Logger
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if session.GetToken(r) != "" {
		http.Redirect(w, r, "/sales", http.StatusSeeOther)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "login-page", nil)
}

func (h *AuthHandler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")



	if email == "" || password == "" {
		w.Header().Set("X-Toast-Message", "Email and password required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.AuthClient.Login(email, password)
	if err != nil {
		h.Log.Error("login failed", "email", email, "error", err)
		w.Header().Set("X-Toast-Message", "Invalid email or password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session.SetToken(w, result.AccessToken)

	// Fetch user's stores and set a default active store
	stores, err := h.AccountClient.GetUserStores(result.User.ID, result.AccessToken)
	if err != nil {
		// Don't block login on this — log it, let them proceed without an active store set.
		// Downstream pages should handle "no active store" gracefully (edge case below).
		h.Log.Error("failed to fetch user stores after login", "user_id", result.User.ID, "error", err)
	} else if len(stores) > 0 {
		session.SetActiveStore(w, stores[0]) // default to first store
	}

	w.Header().Set("HX-Redirect", "/sales")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) SignupPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "signup-page", nil)
}

func (h *AuthHandler) SignupSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("full_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Tmpl.ExecuteTemplate(w, "signup-error", map[string]string{"Error": "All fields are required"})
		return
	}

	// TODO: call auth-service to create account, send OTP
	w.Header().Set("HX-Redirect", "/verify-otp")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) VerifyOTPPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "verify-otp-page", map[string]any{
		"Email": "obi***nenye@gmail.com", // hardcoded, wire from session later
	})
}

func (h *AuthHandler) VerifyOTPSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("otp1") + r.FormValue("otp2") + r.FormValue("otp3") +
		r.FormValue("otp4") + r.FormValue("otp5") + r.FormValue("otp6")

	// TODO: call auth-service to verify
	if len(code) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		h.Tmpl.ExecuteTemplate(w, "verify-otp-error", map[string]string{"Error": "Enter the complete 6-digit code"})
		return
	}

	w.Header().Set("HX-Redirect", "/business-setup")
	w.WriteHeader(http.StatusOK)
}


func (h *AuthHandler) BusinessSetupPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "business-setup-page", nil)
}

func (h* AuthHandler) BusinessSetupSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("business_name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Tmpl.ExecuteTemplate(w, "business-error", map[string]string{"Error": "Enter a business name or tap Skip"})
		return
	}

	// TODO: call business/inventory-service to save business name
	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session.ClearToken(w)
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}