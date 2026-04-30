package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// Invoice represents an SAP A/R Invoice (OINV).
type Invoice struct {
	conn   Executor
	header invoiceHeader
	lines  []invoiceLine
}

type invoiceHeader struct {
	CardCode   string `json:"CardCode,omitempty"`
	DocDate    string `json:"DocDate,omitempty"`
	DocDueDate string `json:"DocDueDate,omitempty"`
	Comments   string `json:"Comments,omitempty"`
}

type invoiceLine struct {
	ItemCode string  `json:"ItemCode,omitempty"`
	Quantity float64 `json:"Quantity,omitempty"`
	Price    float64 `json:"Price,omitempty"`
}

type invoicePayload struct {
	invoiceHeader
	DocumentLines []invoiceLine `json:"DocumentLines"`
}

// NewInvoice initializes a new A/R Invoice object.
func NewInvoice(conn Executor) *Invoice {
	return &Invoice{
		conn: conn,
	}
}

// Header returns a builder to configure the invoice header fields.
func (i *Invoice) Header() *InvoiceHeaderBuilder {
	return &InvoiceHeaderBuilder{i: i}
}

type InvoiceHeaderBuilder struct {
	i *Invoice
}

func (b *InvoiceHeaderBuilder) CardCode(code string) *InvoiceHeaderBuilder {
	b.i.header.CardCode = code
	return b
}

func (b *InvoiceHeaderBuilder) DocDate(date string) *InvoiceHeaderBuilder {
	b.i.header.DocDate = date
	return b
}

func (b *InvoiceHeaderBuilder) DocDueDate(dueDate string) *InvoiceHeaderBuilder {
	b.i.header.DocDueDate = dueDate
	return b
}

func (b *InvoiceHeaderBuilder) Comments(comments string) *InvoiceHeaderBuilder {
	b.i.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single invoice line.
func (i *Invoice) AddLine() *InvoiceLineBuilder {
	return &InvoiceLineBuilder{i: i}
}

type InvoiceLineBuilder struct {
	i    *Invoice
	line invoiceLine
}

func (b *InvoiceLineBuilder) ItemCode(code string) *InvoiceLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *InvoiceLineBuilder) Quantity(qty float64) *InvoiceLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *InvoiceLineBuilder) Price(price float64) *InvoiceLineBuilder {
	b.line.Price = price
	return b
}

func (b *InvoiceLineBuilder) Add() *Invoice {
	b.i.lines = append(b.i.lines, b.line)
	return b.i
}

// Add executes the POST request to create the Invoice in SAP.
func (i *Invoice) Add() (*Response, error) {
	payload := invoicePayload{
		invoiceHeader: i.header,
		DocumentLines: i.lines,
	}

	var rawResp map[string]any
	err := i.conn.Do("POST", "/Invoices", payload, &rawResp)

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
