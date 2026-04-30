package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// Executor represents an interface for executing HTTP requests.
type Executor interface {
	Do(method, path string, reqBody any, resBody any) error
}

// StockTransfer represents a single SAP Stock Transfer document.
type StockTransfer struct {
	conn   Executor
	header StockTransferHeader
	lines  []StockTransferLine
}

// StockTransferHeader contains the header-level fields for a Stock Transfer.
type StockTransferHeader struct {
	FromWarehouse string `json:"FromWarehouse,omitempty"`
	ToWarehouse   string `json:"ToWarehouse,omitempty"`
	Comments      string `json:"Comments,omitempty"`
}

// StockTransferLine contains the line-level fields for a Stock Transfer.
type StockTransferLine struct {
	ItemCode      string  `json:"ItemCode,omitempty"`
	Quantity      float64 `json:"Quantity,omitempty"`
	FromWarehouse string  `json:"FromWarehouse,omitempty"`
	ToWarehouse   string  `json:"ToWarehouse,omitempty"`
}

// stockTransferPayload is the exact JSON structure required by Service Layer.
type stockTransferPayload struct {
	StockTransferHeader
	StockTransferLines []StockTransferLine `json:"StockTransferLines"`
}

// NewStockTransfer initializes a new Stock Transfer object.
func NewStockTransfer(conn Executor) *StockTransfer {
	return &StockTransfer{
		conn: conn,
	}
}

// Header returns a builder to configure the header fields.
func (st *StockTransfer) Header() *StockTransferHeaderBuilder {
	return &StockTransferHeaderBuilder{st: st}
}

// StockTransferHeaderBuilder provides a fluent API for the header.
type StockTransferHeaderBuilder struct {
	st *StockTransfer
}

func (b *StockTransferHeaderBuilder) FromWarehouse(whs string) *StockTransferHeaderBuilder {
	b.st.header.FromWarehouse = whs
	return b
}

func (b *StockTransferHeaderBuilder) ToWarehouse(whs string) *StockTransferHeaderBuilder {
	b.st.header.ToWarehouse = whs
	return b
}

func (b *StockTransferHeaderBuilder) Comments(comments string) *StockTransferHeaderBuilder {
	b.st.header.Comments = comments
	return b
}

// AddLine initializes a new builder for a single stock transfer line.
func (st *StockTransfer) AddLine() *StockTransferLineBuilder {
	return &StockTransferLineBuilder{st: st}
}

// StockTransferLineBuilder provides a fluent API for a document line.
type StockTransferLineBuilder struct {
	st   *StockTransfer
	line StockTransferLine
}

func (b *StockTransferLineBuilder) ItemCode(code string) *StockTransferLineBuilder {
	b.line.ItemCode = code
	return b
}

func (b *StockTransferLineBuilder) Quantity(qty float64) *StockTransferLineBuilder {
	b.line.Quantity = qty
	return b
}

func (b *StockTransferLineBuilder) FromWarehouse(whs string) *StockTransferLineBuilder {
	b.line.FromWarehouse = whs
	return b
}

func (b *StockTransferLineBuilder) ToWarehouse(whs string) *StockTransferLineBuilder {
	b.line.ToWarehouse = whs
	return b
}

// Add appends the constructed line to the main StockTransfer object and returns it for chaining if desired.
func (b *StockTransferLineBuilder) Add() *StockTransfer {
	b.st.lines = append(b.st.lines, b.line)
	return b.st
}

// Add executes the POST request to create the Stock Transfer in SAP.
func (st *StockTransfer) Add() (*Response, error) {
	payload := stockTransferPayload{
		StockTransferHeader: st.header,
		StockTransferLines:  st.lines,
	}

	var rawResp map[string]any
	err := st.conn.Do("POST", "/StockTransfers", payload, &rawResp)

	resp := &Response{}
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()

		// If the error is a structured SAPError, map it to the response object
		if sapErr, ok := err.(*errors.SAPError); ok {
			resp.Error = sapErr
		}
		return resp, err
	}

	resp.Success = true
	resp.Data = rawResp
	return resp, nil
}
