package handlers

import (
	"html/template"

	"github.com/retail-core/xsell-web/internal/clients"
)

type DummyHandler struct {
	Tmpl       *template.Template
	AuthClient *clients.AuthClient
}
