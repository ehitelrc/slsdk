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
	addresses    []bpAddress
	contacts     []contactEmployee
}

type bpAddress struct {
	AddressName string `json:"AddressName,omitempty"`
	Street      string `json:"Street,omitempty"`
	ZipCode     string `json:"ZipCode,omitempty"`
	City        string `json:"City,omitempty"`
	County      string `json:"County,omitempty"`
	Country     string `json:"Country,omitempty"`
	State       string `json:"State,omitempty"`
	AddressType string `json:"AddressType,omitempty"` // "bo_BillTo" or "bo_ShipTo"
}

type contactEmployee struct {
	Name       string `json:"Name,omitempty"`
	Phone1     string `json:"Phone1,omitempty"`
	E_Mail     string `json:"E_Mail,omitempty"`
	Position   string `json:"Position,omitempty"`
	Profession string `json:"Profession,omitempty"`
}

// businessPartnerPayload represents the JSON payload for a Business Partner.
type businessPartnerPayload struct {
	CardCode        string  `json:"CardCode,omitempty"`
	CardName        string  `json:"CardName,omitempty"`
	CardType        string  `json:"CardType,omitempty"`
	FederalTaxID    string  `json:"FederalTaxID,omitempty"`
	GroupCode       int     `json:"GroupCode,omitempty"`
	Currency        string  `json:"Currency,omitempty"`
	Phone1          string  `json:"Phone1,omitempty"`
	Phone2          string  `json:"Phone2,omitempty"`
	Cellular        string  `json:"Cellular,omitempty"`
	Fax             string  `json:"Fax,omitempty"`
	EmailAddress    string  `json:"EmailAddress,omitempty"`
	Notes           string  `json:"Notes,omitempty"`
	DiscountPercent float64 `json:"DiscountPercent,omitempty"`
	CreditLimit     float64 `json:"CreditLimit,omitempty"`
	VatStatus       string  `json:"VatStatus,omitempty"`
	SalesPersonCode int     `json:"SalesPersonCode,omitempty"`
	PayTermsGrpCode int     `json:"PayTermsGrpCode,omitempty"`
}

// NewBusinessPartner initializes a new Business Partner object.
func NewBusinessPartner(conn Executor) *BusinessPartner {
	return &BusinessPartner{
		conn: conn,
	}
}

// --- Header Setters ---

func (bp *BusinessPartner) CardCode(code string) *BusinessPartner {
	bp.payload.CardCode = code
	return bp
}

func (bp *BusinessPartner) CardName(name string) *BusinessPartner {
	bp.payload.CardName = name
	return bp
}

func (bp *BusinessPartner) CardType(cardType string) *BusinessPartner {
	bp.payload.CardType = cardType
	return bp
}

func (bp *BusinessPartner) FederalTaxID(taxID string) *BusinessPartner {
	bp.payload.FederalTaxID = taxID
	return bp
}

func (bp *BusinessPartner) GroupCode(code int) *BusinessPartner {
	bp.payload.GroupCode = code
	return bp
}

func (bp *BusinessPartner) Currency(currency string) *BusinessPartner {
	bp.payload.Currency = currency
	return bp
}

func (bp *BusinessPartner) Phone1(phone string) *BusinessPartner {
	bp.payload.Phone1 = phone
	return bp
}

func (bp *BusinessPartner) Phone2(phone string) *BusinessPartner {
	bp.payload.Phone2 = phone
	return bp
}

func (bp *BusinessPartner) Cellular(cellular string) *BusinessPartner {
	bp.payload.Cellular = cellular
	return bp
}

func (bp *BusinessPartner) Fax(fax string) *BusinessPartner {
	bp.payload.Fax = fax
	return bp
}

func (bp *BusinessPartner) EmailAddress(email string) *BusinessPartner {
	bp.payload.EmailAddress = email
	return bp
}

func (bp *BusinessPartner) Notes(notes string) *BusinessPartner {
	bp.payload.Notes = notes
	return bp
}

func (bp *BusinessPartner) DiscountPercent(discount float64) *BusinessPartner {
	bp.payload.DiscountPercent = discount
	return bp
}

func (bp *BusinessPartner) CreditLimit(limit float64) *BusinessPartner {
	bp.payload.CreditLimit = limit
	return bp
}

func (bp *BusinessPartner) VatStatus(status string) *BusinessPartner {
	bp.payload.VatStatus = status
	return bp
}

func (bp *BusinessPartner) SalesPersonCode(code int) *BusinessPartner {
	bp.payload.SalesPersonCode = code
	return bp
}

func (bp *BusinessPartner) PayTermsGrpCode(code int) *BusinessPartner {
	bp.payload.PayTermsGrpCode = code
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

// --- BP Addresses ---

func (bp *BusinessPartner) AddAddress() *BPAddressBuilder {
	return &BPAddressBuilder{bp: bp}
}

type BPAddressBuilder struct {
	bp   *BusinessPartner
	addr bpAddress
}

func (b *BPAddressBuilder) AddressName(name string) *BPAddressBuilder {
	b.addr.AddressName = name
	return b
}

func (b *BPAddressBuilder) Street(street string) *BPAddressBuilder {
	b.addr.Street = street
	return b
}

func (b *BPAddressBuilder) ZipCode(zip string) *BPAddressBuilder {
	b.addr.ZipCode = zip
	return b
}

func (b *BPAddressBuilder) City(city string) *BPAddressBuilder {
	b.addr.City = city
	return b
}

func (b *BPAddressBuilder) County(county string) *BPAddressBuilder {
	b.addr.County = county
	return b
}

func (b *BPAddressBuilder) Country(country string) *BPAddressBuilder {
	b.addr.Country = country
	return b
}

func (b *BPAddressBuilder) State(state string) *BPAddressBuilder {
	b.addr.State = state
	return b
}

func (b *BPAddressBuilder) AddressType(typeStr string) *BPAddressBuilder {
	b.addr.AddressType = typeStr
	return b
}

func (b *BPAddressBuilder) Add() *BusinessPartner {
	b.bp.addresses = append(b.bp.addresses, b.addr)
	return b.bp
}

// --- Contact Employees ---

func (bp *BusinessPartner) AddContact() *BPContactBuilder {
	return &BPContactBuilder{bp: bp}
}

type BPContactBuilder struct {
	bp   *BusinessPartner
	cont contactEmployee
}

func (b *BPContactBuilder) Name(name string) *BPContactBuilder {
	b.cont.Name = name
	return b
}

func (b *BPContactBuilder) Phone1(phone string) *BPContactBuilder {
	b.cont.Phone1 = phone
	return b
}

func (b *BPContactBuilder) E_Mail(email string) *BPContactBuilder {
	b.cont.E_Mail = email
	return b
}

func (b *BPContactBuilder) Position(position string) *BPContactBuilder {
	b.cont.Position = position
	return b
}

func (b *BPContactBuilder) Profession(profession string) *BPContactBuilder {
	b.cont.Profession = profession
	return b
}

func (b *BPContactBuilder) Add() *BusinessPartner {
	b.bp.contacts = append(b.bp.contacts, b.cont)
	return b.bp
}

// --- Execution ---

// Add executes the POST request to create the Business Partner in SAP.
func (bp *BusinessPartner) Add() (*Response, error) {
	// Merge static payload with custom fields and sub-collections
	var merged map[string]any

	payloadBytes, _ := json.Marshal(bp.payload)
	json.Unmarshal(payloadBytes, &merged)

	if len(bp.addresses) > 0 {
		merged["BPAddresses"] = bp.addresses
	}
	if len(bp.contacts) > 0 {
		merged["ContactEmployees"] = bp.contacts
	}

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
