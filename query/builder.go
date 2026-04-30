package query

import (
	"fmt"
	"net/url"
	"strings"
)

// Executor represents an interface for executing HTTP requests.
type Executor interface {
	Do(method, path string, reqBody any, resBody any) error
}

// Builder provides a fluent API for constructing OData queries.
type Builder struct {
	conn     Executor
	resource string
	selects  []string
	filters  []string
	orderBys []string
	top      int
	skip     int
}

// NewBuilder initializes a new OData query builder.
func NewBuilder(resource string, conn Executor) *Builder {
	return &Builder{
		conn:     conn,
		resource: resource,
	}
}

// Select specifies the fields to return.
func (b *Builder) Select(fields ...string) *Builder {
	b.selects = append(b.selects, fields...)
	return b
}

// Filter adds a filter condition.
func (b *Builder) Filter(filter string) *Builder {
	b.filters = append(b.filters, filter)
	return b
}

// OrderBy specifies the sorting order.
func (b *Builder) OrderBy(orderBy string) *Builder {
	b.orderBys = append(b.orderBys, orderBy)
	return b
}

// Top limits the number of returned records.
func (b *Builder) Top(top int) *Builder {
	b.top = top
	return b
}

// Skip bypasses a specified number of records.
func (b *Builder) Skip(skip int) *Builder {
	b.skip = skip
	return b
}

// Get executes the built OData query against the Service Layer.
func (b *Builder) Get() (any, error) {
	q := url.Values{}

	if len(b.selects) > 0 {
		q.Set("$select", strings.Join(b.selects, ","))
	}
	if len(b.filters) > 0 {
		q.Set("$filter", strings.Join(b.filters, " and "))
	}
	if len(b.orderBys) > 0 {
		q.Set("$orderby", strings.Join(b.orderBys, ","))
	}
	if b.top > 0 {
		q.Set("$top", fmt.Sprintf("%d", b.top))
	}
	if b.skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", b.skip))
	}

	path := fmt.Sprintf("/%s", b.resource)
	if len(q) > 0 {
		path = fmt.Sprintf("%s?%s", path, q.Encode())
	}

	var result map[string]any
	err := b.conn.Do("GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	// For OData, collections are typically in the "value" property
	if val, ok := result["value"]; ok {
		return val, nil
	}

	return result, nil
}
