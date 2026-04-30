package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// Order represents an SAP Sales Order (ORDR).
type Order struct {
	conn   Executor
	header orderHeader
	lines  []orderLine
}

type orderHeader struct {
	CardCode   string `json:"CardCode,omitempty"`
	DocDate    string `json:"DocDate,omitempty"`
	DocDueDate string `json:"DocDueDate,omitempty"`
	Comments   string `json:"Comments,omitempty"`
}

type orderLine struct {
	ItemCode string  `json:"ItemCode,omitempty"`
	Quantity float64 `json:"Quantity,omitempty"`
	Price    float64 `json:"Price,omitempty"`
}

type orderPayload struct {
	orderHeader
	DocumentLines []orderLine `json:"DocumentLines"`
}

// NewOrder initializes a new Sales Order object.
func NewOrder(conn Executor) *Order {
	return &Order{
		conn: conn,
	}
}

// Header returns a builder to configure the order header fields.
func (o *Order) Header() *OrderHeaderBuilder {
	return &OrderHeaderBuilder{o: o}
}

type OrderHeaderBuilder struct {
	o *Order
}

func (b *OrderHeaderBuilder) CardCode(code string) *OrderHeaderBuilder {
	b.o.header.CardCode = code
	return b
}

func (b *OrderHeaderBuilder) DocDate(date string) *OrderHeaderBuilder {
	b.o.header.DocDate = date
	return b
}

func (b *OrderHeaderBuilder) DocDueDate(dueDate string) *OrderHeaderBuilder {
	b.o.header.DocDueDate = dueDate
	return b
}

func (b *OrderHeaderBuilder) Comments(comments string) *OrderHeaderBuilder {
	b.o.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single order line.
func (o *Order) AddLine() *OrderLineBuilder {
	return &OrderLineBuilder{o: o}
}

type OrderLineBuilder struct {
	o    *Order
	line orderLine
}

func (b *OrderLineBuilder) ItemCode(code string) *OrderLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *OrderLineBuilder) Quantity(qty float64) *OrderLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *OrderLineBuilder) Price(price float64) *OrderLineBuilder {
	b.line.Price = price
	return b
}

func (b *OrderLineBuilder) Add() *Order {
	b.o.lines = append(b.o.lines, b.line)
	return b.o
}

// Add executes the POST request to create the Sales Order in SAP.
func (o *Order) Add() (*Response, error) {
	payload := orderPayload{
		orderHeader:   o.header,
		DocumentLines: o.lines,
	}

	var rawResp map[string]any
	err := o.conn.Do("POST", "/Orders", payload, &rawResp)

	resp := &Response{}
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()

		if sapErr, ok := err.(*errors.SAPError); ok {
			resp.Error = sapErr
		}
		return resp, err
	}

	resp.Success = true
	resp.Data = rawResp
	return resp, nil
}
