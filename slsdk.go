package slsdk

import (
	"github.com/ehitelrc/slsdk/connection"
	"github.com/ehitelrc/slsdk/objects"
)

// Config represents the credentials and settings required to connect to the Service Layer.
type Config = connection.Config

// Connection manages the HTTP client and handles the session.
type Connection = connection.Connection

// NewConnection initializes a new Service Layer connection manager.
func NewConnection(cfg Config) *Connection {
	return connection.NewConnection(cfg)
}

// NewStockTransfer initializes a new Stock Transfer object.
func NewStockTransfer(conn *Connection) *objects.StockTransfer {
	return objects.NewStockTransfer(conn)
}

// NewBusinessPartner initializes a new Business Partner object.
func NewBusinessPartner(conn *Connection) *objects.BusinessPartner {
	return objects.NewBusinessPartner(conn)
}

// NewItem initializes a new Item object.
func NewItem(conn *Connection) *objects.Item {
	return objects.NewItem(conn)
}

// NewOrder initializes a new Sales Order object.
func NewOrder(conn *Connection) *objects.Order {
	return objects.NewOrder(conn)
}

// NewInvoice initializes a new A/R Invoice object.
func NewInvoice(conn *Connection) *objects.Invoice {
	return objects.NewInvoice(conn)
}

// NewPurchaseOrder initializes a new Purchase Order object.
func NewPurchaseOrder(conn *Connection) *objects.PurchaseOrder {
	return objects.NewPurchaseOrder(conn)
}

