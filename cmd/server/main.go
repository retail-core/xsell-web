package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/retail-core/xsell-web/internal/clients"
	"github.com/retail-core/xsell-web/internal/handlers"
	"github.com/retail-core/xsell-web/internal/logger"
	m_middleware "github.com/retail-core/xsell-web/internal/middleware"
	"github.com/retail-core/xsell-web/internal/templates"
)

func main() {
	__logger := logger.New()

	if err := godotenv.Load(); err != nil {
		__logger.Warn("Warning: .env file not found:", "err", err)
	}

	tmpl, err := templates.Load()
	if err != nil {
		log.Fatal("failed to load templates:", err)
	}

	authClient := clients.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
	accountClient := clients.NewAccountClient(os.Getenv("ACCOUNT_SERVICE_URL"))

	authHandler := &handlers.AuthHandler{
		Tmpl:          tmpl,
		AuthClient:    authClient,
		AccountClient: accountClient,
		Log:           __logger,
	}

	accountHandler := &handlers.AccountHandler{
		ShellHandler:  handlers.ShellHandler{Tmpl: tmpl, Log: __logger},
		AccountClient: accountClient,
	}

	salesClient := clients.NewSalesClient(os.Getenv("SALES_SERVICE_URL"))
	salesHandler := &handlers.SalesHandler{
		ShellHandler: handlers.ShellHandler{Tmpl: tmpl, Log: __logger},
		SalesClient:  salesClient,
	}

	inventoryClient := clients.NewInventoryClient(os.Getenv("INVENTORY_SERVICE_URL"))
	inventoryHandler := &handlers.InventoryHandler{
		ShellHandler:    handlers.ShellHandler{Tmpl: tmpl, Log: __logger},
		InventoryClient: inventoryClient,
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(m_middleware.Recover(__logger))

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// auth
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.LoginSubmit)

	r.Get("/signup", authHandler.SignupPage)
	r.Post("/signup", authHandler.SignupSubmit)
	r.Get("/verify-otp", authHandler.VerifyOTPPage)
	r.Post("/verify-otp", authHandler.VerifyOTPSubmit)
	r.Get("/business-setup", authHandler.BusinessSetupPage)
	r.Post("/business-setup", authHandler.BusinessSetupSubmit)
	r.Post("/logout", authHandler.Logout)

	r.Get("/access-denied", accountHandler.AccessDeniedPage) // public, no guard needed

	// main
	r.Group(func(protected chi.Router) {
		protected.Use(m_middleware.RequireAuth)

		protected.Get("/products", inventoryHandler.Products)
		protected.Get("/pos", inventoryHandler.POS)
		protected.Get("/sales", salesHandler.Sales)
		protected.Get("/account", accountHandler.Account)
		protected.Get("/stores", accountHandler.StoresPage)
		protected.Get("/settings", accountHandler.SettingsPage)
		protected.Get("/store-switcher", accountHandler.StoreSwitcherSheet)
		protected.Post("/switch-store", accountHandler.SwitchStore)
		protected.Post("/pos/checkout", salesHandler.Checkout)
		protected.Get("/orders/{id}", salesHandler.OrderDetail)
	})

	r.Group(func(owner chi.Router) {
		owner.Use(m_middleware.RequireAuth)
		owner.Use(m_middleware.RequireOwner)

		owner.Get("/staff", accountHandler.StaffPage)
		owner.Get("/staff/new", accountHandler.StaffFormPage)
		owner.Get("/staff/{id}/edit", accountHandler.StaffFormPage)
		owner.Get("/stores", accountHandler.StoresPage)
		owner.Get("/stores/new", accountHandler.StoreFormPage)
		owner.Get("/stores/{id}/edit", accountHandler.StoreFormPage)
		owner.Get("/subscription", accountHandler.SubscriptionPage)
		owner.Get("/products/new", inventoryHandler.ProductFormPage)
		owner.Get("/products/{id}/edit", inventoryHandler.ProductFormPage)
		owner.Post("/products/{id}/restock", inventoryHandler.RestockSubmit)
		owner.Delete("/products/{id}", inventoryHandler.DeleteProduct)
		owner.Post("/products/new", inventoryHandler.ProductFormSubmit)
		owner.Post("/products/{id}/edit", inventoryHandler.ProductFormSubmit)
		owner.Patch("/products/{id}/status", inventoryHandler.ToggleStatus)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		tmpl.ExecuteTemplate(w, "not-found-page", nil)
	})

	__logger.Info("xsell-web listening on :8000")
	http.ListenAndServe(":8000", r)
}
