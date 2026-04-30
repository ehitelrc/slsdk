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
