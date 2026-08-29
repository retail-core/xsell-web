package handlers

import (
	"github.com/retail-core/xsell-web/internal/clients"
)

type InventoryHandler struct {
	ShellHandler
	InventoryClient *clients.InventoryClient
}