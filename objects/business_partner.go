package objects

import (
	"encoding/json"

	"github.com/ehitelrc/slsdk/errors"
)

// BusinessPartner represents a single SAP Business Partner.
type BusinessPartner struct {
	conn         Executor
	payload      businessPartnerPayload
	customFields map[string]any
}

// businessPartnerPayload represents the JSON payload for a Business Partner.
type businessPartnerPayload struct {
	CardCode     string `json:"CardCode,omitempty"`
	CardName     string `json:"CardName,omitempty"`
	CardType     string `json:"CardType,omitempty"`
	FederalTaxID string `json:"FederalTaxID,omitempty"`
}

// NewBusinessPartner initializes a new Business Partner object.
func NewBusinessPartner(conn Executor) *BusinessPartner {
	return &BusinessPartner{
		conn: conn,
	}
}

// CardCode sets the Business Partner Code.
func (bp *BusinessPartner) CardCode(code string) *BusinessPartner {
	bp.payload.CardCode = code
	return bp
}

// CardName sets the Business Partner Name.
func (bp *BusinessPartner) CardName(name string) *BusinessPartner {
	bp.payload.CardName = name
	return bp
}

// CardType sets the Business Partner Type (e.g., "C" for Customer, "S" for Supplier, "L" for Lead).
func (bp *BusinessPartner) CardType(cardType string) *BusinessPartner {
	bp.payload.CardType = cardType
	return bp
}

// FederalTaxID sets the Federal Tax ID (RFC in Mexico localization).
func (bp *BusinessPartner) FederalTaxID(taxID string) *BusinessPartner {
	bp.payload.FederalTaxID = taxID
	return bp
}

// Set allows setting any arbitrary field or User Defined Field (UDF) by its string name.
func (bp *BusinessPartner) Set(field string, value any) *BusinessPartner {
	if bp.customFields == nil {
		bp.customFields = make(map[string]any)
	}
	bp.customFields[field] = value
	return bp
}

// Add executes the POST request to create the Business Partner in SAP.
func (bp *BusinessPartner) Add() (*Response, error) {
	// Merge static payload with custom fields
	var merged map[string]any
	
	// Convert static payload to map using JSON to respect 'omitempty'
	payloadBytes, _ := json.Marshal(bp.payload)
	json.Unmarshal(payloadBytes, &merged)

	// Inject custom fields
	for k, v := range bp.customFields {
		merged[k] = v
	}

	var rawResp map[string]any
	err := bp.conn.Do("POST", "/BusinessPartners", merged, &rawResp)

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
