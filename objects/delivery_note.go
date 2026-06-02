package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// DeliveryNote represents an SAP Delivery Note (ODLN - AR Delivery).
type DeliveryNote struct {
	conn   Executor
	header dnHeader
	lines  []dnLine
}

type dnHeader struct {
	CardCode string `json:"CardCode,omitempty"`
	DocDate  string `json:"DocDate,omitempty"`
	Comments string `json:"Comments,omitempty"`
}

type dnLine struct {
	BaseType      int        `json:"BaseType"`
	BaseEntry     int        `json:"BaseEntry"`
	BaseLine      int        `json:"BaseLine"`
	ItemCode      string     `json:"ItemCode,omitempty"`
	Quantity      float64    `json:"Quantity,omitempty"`
	UseBaseUnits  string     `json:"UseBaseUnits,omitempty"`
	WarehouseCode string     `json:"WarehouseCode,omitempty"`
	BatchNumbers  []dnBatch  `json:"BatchNumbers,omitempty"`
	SerialNumbers []dnSerial `json:"SerialNumbers,omitempty"`
}

type dnBatch struct {
	BatchNumber string  `json:"BatchNumber"`
	Quantity    float64 `json:"Quantity"`
}

type dnSerial struct {
	InternalSerialNumber string `json:"InternalSerialNumber"`
	SystemSerialNumber   int    `json:"SystemSerialNumber,omitempty"`
}

type dnPayload struct {
	dnHeader
	DocumentLines []dnLine `json:"DocumentLines"`
}

// NewDeliveryNote initializes a new Delivery Note object.
func NewDeliveryNote(conn Executor) *DeliveryNote {
	return &DeliveryNote{
		conn: conn,
	}
}

// Header returns a builder to configure the header fields.
func (p *DeliveryNote) Header() *DNHeaderBuilder {
	return &DNHeaderBuilder{p: p}
}

type DNHeaderBuilder struct {
	p *DeliveryNote
}

func (b *DNHeaderBuilder) CardCode(code string) *DNHeaderBuilder {
	b.p.header.CardCode = code
	return b
}

func (b *DNHeaderBuilder) DocDate(date string) *DNHeaderBuilder {
	b.p.header.DocDate = date
	return b
}

func (b *DNHeaderBuilder) Comments(comments string) *DNHeaderBuilder {
	b.p.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single delivery note line.
func (p *DeliveryNote) AddLine() *DNLineBuilder {
	return &DNLineBuilder{
		p: p,
		line: dnLine{
			BatchNumbers:  make([]dnBatch, 0),
			SerialNumbers: make([]dnSerial, 0),
		},
	}
}

type DNLineBuilder struct {
	p    *DeliveryNote
	line dnLine
}

func (b *DNLineBuilder) BaseType(baseType int) *DNLineBuilder {
	b.line.BaseType = baseType
	return b
}

func (b *DNLineBuilder) BaseEntry(baseEntry int) *DNLineBuilder {
	b.line.BaseEntry = baseEntry
	return b
}

func (b *DNLineBuilder) BaseLine(baseLine int) *DNLineBuilder {
	b.line.BaseLine = baseLine
	return b
}

func (b *DNLineBuilder) ItemCode(code string) *DNLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *DNLineBuilder) Quantity(qty float64) *DNLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *DNLineBuilder) UseBaseUnits(val string) *DNLineBuilder {
	b.line.UseBaseUnits = val
	return b
}

func (b *DNLineBuilder) WarehouseCode(code string) *DNLineBuilder {
	b.line.WarehouseCode = code
	return b
}

// AddBatch adds a batch assignment to the current line.
func (b *DNLineBuilder) AddBatch(batchNumber string, qty float64) *DNLineBuilder {
	b.line.BatchNumbers = append(b.line.BatchNumbers, dnBatch{
		BatchNumber: batchNumber,
		Quantity:    qty,
	})
	return b
}

// AddSerial adds a serial number assignment to the current line.
func (b *DNLineBuilder) AddSerial(serialNumber string) *DNLineBuilder {
	b.line.SerialNumbers = append(b.line.SerialNumbers, dnSerial{
		InternalSerialNumber: serialNumber,
	})
	return b
}

// Add appends the constructed line to the document and returns the main builder.
func (b *DNLineBuilder) Add() *DeliveryNote {
	b.p.lines = append(b.p.lines, b.line)
	return b.p
}

// GetPayload returns the payload structure for debugging.
func (p *DeliveryNote) GetPayload() interface{} {
	return dnPayload{
		dnHeader:      p.header,
		DocumentLines: p.lines,
	}
}

// Add executes the POST request to create the Delivery Note in SAP.
func (p *DeliveryNote) Add() (*Response, error) {
	payload := p.GetPayload()

	var rawResp map[string]any
	err := p.conn.Do("POST", "/DeliveryNotes", payload, &rawResp)

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
