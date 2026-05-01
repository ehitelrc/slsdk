package objects

import (
	"github.com/ehitelrc/slsdk/errors"
)

// Map is a convenient alias for map[string]any to use in generic payloads.
type Map map[string]any

// GenericObject allows interacting with unmapped Service Layer endpoints using a dynamic payload.
type GenericObject struct {
	conn       Executor
	endpoint   string
	properties Map
}

// NewGenericObject initializes a new generic object targeted at a specific endpoint.
func NewGenericObject(conn Executor, endpoint string) *GenericObject {
	return &GenericObject{
		conn:       conn,
		endpoint:   endpoint,
		properties: make(Map),
	}
}

// Set assigns a value to a field in the generic payload.
func (g *GenericObject) Set(field string, value any) *GenericObject {
	g.properties[field] = value
	return g
}

// Append adds an item (e.g., a line or sub-object) to an array collection within the payload.
func (g *GenericObject) Append(collectionName string, item any) *GenericObject {
	if _, ok := g.properties[collectionName]; !ok {
		g.properties[collectionName] = []any{}
	}

	collection := g.properties[collectionName].([]any)
	g.properties[collectionName] = append(collection, item)
	return g
}

// Add executes a POST request to create the object in SAP.
func (g *GenericObject) Add() (*Response, error) {
	var rawResp map[string]any
	err := g.conn.Do("POST", "/"+g.endpoint, g.properties, &rawResp)

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
