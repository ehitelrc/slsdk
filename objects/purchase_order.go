package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// PurchaseOrder represents an SAP Purchase Order (OPOR).
type PurchaseOrder struct {
	conn   Executor
	header poHeader
	lines  []poLine
}

type poHeader struct {
	CardCode   string `json:"CardCode,omitempty"`
	DocDate    string `json:"DocDate,omitempty"`
	DocDueDate string `json:"DocDueDate,omitempty"`
	Comments   string `json:"Comments,omitempty"`
}

type poLine struct {
	ItemCode string  `json:"ItemCode,omitempty"`
	Quantity float64 `json:"Quantity,omitempty"`
	Price    float64 `json:"Price,omitempty"`
}

type poPayload struct {
	poHeader
	DocumentLines []poLine `json:"DocumentLines"`
}

// NewPurchaseOrder initializes a new Purchase Order object.
func NewPurchaseOrder(conn Executor) *PurchaseOrder {
	return &PurchaseOrder{
		conn: conn,
	}
}

// Header returns a builder to configure the Purchase Order header fields.
func (p *PurchaseOrder) Header() *POHeaderBuilder {
	return &POHeaderBuilder{p: p}
}

type POHeaderBuilder struct {
	p *PurchaseOrder
}

func (b *POHeaderBuilder) CardCode(code string) *POHeaderBuilder {
	b.p.header.CardCode = code
	return b
}

func (b *POHeaderBuilder) DocDate(date string) *POHeaderBuilder {
	b.p.header.DocDate = date
	return b
}

func (b *POHeaderBuilder) DocDueDate(dueDate string) *POHeaderBuilder {
	b.p.header.DocDueDate = dueDate
	return b
}

func (b *POHeaderBuilder) Comments(comments string) *POHeaderBuilder {
	b.p.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single purchase order line.
func (p *PurchaseOrder) AddLine() *POLineBuilder {
	return &POLineBuilder{p: p}
}

type POLineBuilder struct {
	p    *PurchaseOrder
	line poLine
}

func (b *POLineBuilder) ItemCode(code string) *POLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *POLineBuilder) Quantity(qty float64) *POLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *POLineBuilder) Price(price float64) *POLineBuilder {
	b.line.Price = price
	return b
}

func (b *POLineBuilder) Add() *PurchaseOrder {
	b.p.lines = append(b.p.lines, b.line)
	return b.p
}

// Add executes the POST request to create the Purchase Order in SAP.
func (p *PurchaseOrder) Add() (*Response, error) {
	payload := poPayload{
		poHeader:      p.header,
		DocumentLines: p.lines,
	}

	var rawResp map[string]any
	err := p.conn.Do("POST", "/PurchaseOrders", payload, &rawResp)

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
