package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// Item represents an SAP Master Data Item.
type Item struct {
	conn    Executor
	payload itemPayload
}

type itemPayload struct {
	ItemCode       string `json:"ItemCode,omitempty"`
	ItemName       string `json:"ItemName,omitempty"`
	ItemType       string `json:"ItemType,omitempty"`
	ItemsGroupCode int    `json:"ItemsGroupCode,omitempty"`
}

// NewItem initializes a new Item object.
func NewItem(conn Executor) *Item {
	return &Item{
		conn: conn,
	}
}

// ItemCode sets the Item Code.
func (i *Item) ItemCode(code string) *Item {
	i.payload.ItemCode = code
	return i
}

// ItemName sets the Item Name.
func (i *Item) ItemName(name string) *Item {
	i.payload.ItemName = name
	return i
}

// ItemType sets the type of the item (e.g., "itItems").
func (i *Item) ItemType(itemType string) *Item {
	i.payload.ItemType = itemType
	return i
}

// ItemsGroupCode sets the Group Code for the Item.
func (i *Item) ItemsGroupCode(groupCode int) *Item {
	i.payload.ItemsGroupCode = groupCode
	return i
}

// Add executes the POST request to create the Item in SAP.
func (i *Item) Add() (*Response, error) {
	var rawResp map[string]any
	err := i.conn.Do("POST", "/Items", i.payload, &rawResp)

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
