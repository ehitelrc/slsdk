package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// PurchaseDelivery represents an SAP Purchase Delivery Note (OPDN - Goods Receipt PO).
type PurchaseDelivery struct {
	conn   Executor
	header pdHeader
	lines  []pdLine
}

type pdHeader struct {
	CardCode string `json:"CardCode,omitempty"`
	DocDate  string `json:"DocDate,omitempty"`
	Comments string `json:"Comments,omitempty"`
}

type pdLine struct {
	BaseType      int        `json:"BaseType"`
	BaseEntry     int        `json:"BaseEntry"`
	BaseLine      int        `json:"BaseLine"`
	ItemCode      string     `json:"ItemCode,omitempty"`
	Quantity      float64    `json:"Quantity,omitempty"`
	UseBaseUnits  string     `json:"UseBaseUnits,omitempty"`
	WarehouseCode string     `json:"WarehouseCode,omitempty"`
	BatchNumbers  []pdBatch  `json:"BatchNumbers,omitempty"`
	SerialNumbers []pdSerial `json:"SerialNumbers,omitempty"`
}

type pdBatch struct {
	BatchNumber string  `json:"BatchNumberProperty"`
	Quantity    float64 `json:"Quantity"`
	ExpiryDate  string  `json:"ExpiryDate,omitempty"`
}

type pdSerial struct {
	InternalSerialNumber string `json:"InternalSerialNumber"`
}

type pdPayload struct {
	pdHeader
	DocumentLines []pdLine `json:"DocumentLines"`
}

// NewPurchaseDelivery initializes a new Purchase Delivery Note object.
func NewPurchaseDelivery(conn Executor) *PurchaseDelivery {
	return &PurchaseDelivery{
		conn: conn,
	}
}

// Header returns a builder to configure the header fields.
func (p *PurchaseDelivery) Header() *PDHeaderBuilder {
	return &PDHeaderBuilder{p: p}
}

type PDHeaderBuilder struct {
	p *PurchaseDelivery
}

func (b *PDHeaderBuilder) CardCode(code string) *PDHeaderBuilder {
	b.p.header.CardCode = code
	return b
}

func (b *PDHeaderBuilder) DocDate(date string) *PDHeaderBuilder {
	b.p.header.DocDate = date
	return b
}

func (b *PDHeaderBuilder) Comments(comments string) *PDHeaderBuilder {
	b.p.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single delivery note line.
func (p *PurchaseDelivery) AddLine() *PDLineBuilder {
	return &PDLineBuilder{
		p: p,
		line: pdLine{
			BatchNumbers:  make([]pdBatch, 0),
			SerialNumbers: make([]pdSerial, 0),
		},
	}
}

type PDLineBuilder struct {
	p    *PurchaseDelivery
	line pdLine
}

func (b *PDLineBuilder) BaseType(baseType int) *PDLineBuilder {
	b.line.BaseType = baseType
	return b
}

func (b *PDLineBuilder) BaseEntry(baseEntry int) *PDLineBuilder {
	b.line.BaseEntry = baseEntry
	return b
}

func (b *PDLineBuilder) BaseLine(baseLine int) *PDLineBuilder {
	b.line.BaseLine = baseLine
	return b
}

func (b *PDLineBuilder) ItemCode(code string) *PDLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *PDLineBuilder) Quantity(qty float64) *PDLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *PDLineBuilder) UseBaseUnits(val string) *PDLineBuilder {
	b.line.UseBaseUnits = val
	return b
}

func (b *PDLineBuilder) WarehouseCode(code string) *PDLineBuilder {
	b.line.WarehouseCode = code
	return b
}

// AddBatch adds a batch assignment to the current line.
func (b *PDLineBuilder) AddBatch(batchNumber string, qty float64, expiryDate string) *PDLineBuilder {
	b.line.BatchNumbers = append(b.line.BatchNumbers, pdBatch{
		BatchNumber: batchNumber,
		Quantity:    qty,
		ExpiryDate:  expiryDate,
	})
	return b
}

// AddSerial adds a serial number assignment to the current line.
func (b *PDLineBuilder) AddSerial(serialNumber string) *PDLineBuilder {
	b.line.SerialNumbers = append(b.line.SerialNumbers, pdSerial{
		InternalSerialNumber: serialNumber,
	})
	return b
}

// Add appends the constructed line to the document and returns the main builder.
func (b *PDLineBuilder) Add() *PurchaseDelivery {
	b.p.lines = append(b.p.lines, b.line)
	return b.p
}

// GetPayload returns the payload structure for debugging.
func (p *PurchaseDelivery) GetPayload() interface{} {
	return pdPayload{
		pdHeader:      p.header,
		DocumentLines: p.lines,
	}
}

// Add executes the POST request to create the Purchase Delivery Note in SAP.
func (p *PurchaseDelivery) Add() (*Response, error) {
	payload := p.GetPayload()

	var rawResp map[string]any
	err := p.conn.Do("POST", "/PurchaseDeliveryNotes", payload, &rawResp)

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
