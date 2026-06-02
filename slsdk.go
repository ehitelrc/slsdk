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

// NewPurchaseDelivery initializes a new Purchase Delivery Note (Goods Receipt PO).
func NewPurchaseDelivery(conn *Connection) *objects.PurchaseDelivery {
	return objects.NewPurchaseDelivery(conn)
}

// NewDeliveryNote initializes a new Delivery Note (AR Delivery).
func NewDeliveryNote(conn *Connection) *objects.DeliveryNote {
	return objects.NewDeliveryNote(conn)
}

// Map is a convenient alias for map[string]any to use in generic payloads.
type Map = objects.Map

// NewGenericObject initializes a dynamic object for an unmapped endpoint.
func NewGenericObject(conn *Connection, endpoint string) *objects.GenericObject {
	return objects.NewGenericObject(conn, endpoint)
}

// NewQuery initializes a generic OData query builder for the specified endpoint.
func NewQuery(conn *Connection, endpoint string) *objects.QueryBuilder {
	return objects.NewQuery(conn, endpoint)
}

// NewSQLView initializes an object to interact with an exposed MSSQL View.
func NewSQLView(conn *Connection, viewName string) *objects.SQLView {
	return objects.NewSQLView(conn, viewName)
}

// NewCrossJoin initializes a cross join query targeting QueryService_PostQuery.
func NewCrossJoin(conn *Connection, tables ...string) *objects.CrossJoinBuilder {
	return objects.NewCrossJoin(conn, tables...)
}

