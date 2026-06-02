package objects

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ehitelrc/slsdk/errors"
)

// QueryBuilder provides a fluent interface for building OData queries.
type QueryBuilder struct {
	conn     Executor
	endpoint string
	selects  []string
	filters  []string
	orderBys []string
	top      int
	skip     int
}

// NewQuery initializes a new query for a specific endpoint.
func NewQuery(conn Executor, endpoint string) *QueryBuilder {
	return &QueryBuilder{
		conn:     conn,
		endpoint: endpoint,
	}
}

// Select specifies the fields to retrieve.
func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	q.selects = append(q.selects, fields...)
	return q
}

// Filter specifies the OData filter condition.
func (q *QueryBuilder) Filter(filter string) *QueryBuilder {
	q.filters = append(q.filters, filter)
	return q
}

// OrderBy specifies the sorting order.
func (q *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	q.orderBys = append(q.orderBys, orderBy)
	return q
}

// Top specifies the maximum number of records to return.
func (q *QueryBuilder) Top(top int) *QueryBuilder {
	q.top = top
	return q
}

// Skip specifies the number of records to skip.
func (q *QueryBuilder) Skip(skip int) *QueryBuilder {
	q.skip = skip
	return q
}

// Get executes the GET request with the built OData query options.
func (q *QueryBuilder) Get() (*Response, error) {
	queryVals := url.Values{}

	if len(q.selects) > 0 {
		queryVals.Add("$select", strings.Join(q.selects, ","))
	}
	if len(q.filters) > 0 {
		queryVals.Add("$filter", strings.Join(q.filters, " and "))
	}
	if len(q.orderBys) > 0 {
		queryVals.Add("$orderby", strings.Join(q.orderBys, ","))
	}
	if q.top > 0 {
		queryVals.Add("$top", fmt.Sprintf("%d", q.top))
	}
	if q.skip > 0 {
		queryVals.Add("$skip", fmt.Sprintf("%d", q.skip))
	}

	path := "/" + q.endpoint
	if queryString := queryVals.Encode(); queryString != "" {
		// url.Values encodes spaces as +, but OData often expects %20
		path = path + "?" + strings.ReplaceAll(queryString, "+", "%20")
	}

	var rawResp map[string]any
	err := q.conn.Do("GET", path, nil, &rawResp)

	return buildResponse(rawResp, err)
}

// GetAll automatically handles server-driven pagination by fetching all pages iteratively.
// It executes the initial query, inspects for a NextLink, and continues fetching until all records are retrieved.
// The returned Response will contain a combined "value" array in Data.
func (q *QueryBuilder) GetAll() (*Response, error) {
	// Execute the first request
	resp, err := q.Get()
	if err != nil {
		return nil, err
	}

	// Extract the first page of values
	var allValues []any
	if dataMap, ok := resp.Data.(map[string]any); ok {
		if values, ok := dataMap["value"].([]any); ok {
			allValues = append(allValues, values...)
		}
	}

	// Iterate through NextLinks if they exist
	nextLink := resp.NextLink
	for nextLink != "" {
		var nextResp map[string]any
		
		path := nextLink
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		err := q.conn.Do("GET", path, nil, &nextResp)
		if err != nil {
			return nil, err
		}

		pageResp, pageErr := buildResponse(nextResp, err)
		if pageErr != nil {
			return nil, pageErr
		}

		if dataMap, ok := pageResp.Data.(map[string]any); ok {
			if values, ok := dataMap["value"].([]any); ok {
				allValues = append(allValues, values...)
			}
		}

		nextLink = pageResp.NextLink
	}

	// Overwrite the final response's value with the aggregated array
	if finalData, ok := resp.Data.(map[string]any); ok {
		finalData["value"] = allValues
		resp.Data = finalData
		resp.NextLink = "" // Reached the end
	}

	return resp, nil
}

// SQLView provides methods to interact with MSSQL Views exposed in Service Layer.
type SQLView struct {
	conn     Executor
	viewName string
}

// NewSQLView initializes a SQLView object.
func NewSQLView(conn Executor, viewName string) *SQLView {
	return &SQLView{
		conn:     conn,
		viewName: viewName,
	}
}

// Expose registers the SQL View in the Service Layer.
func (s *SQLView) Expose() (*Response, error) {
	path := "/SQLViews('" + s.viewName + "')/Expose"
	var rawResp map[string]any
	err := s.conn.Do("POST", path, nil, &rawResp)
	return buildResponse(rawResp, err)
}

// Unexpose removes the SQL View from the Service Layer.
func (s *SQLView) Unexpose() (*Response, error) {
	path := "/SQLViews('" + s.viewName + "')/Unexpose"
	var rawResp map[string]any
	err := s.conn.Do("POST", path, nil, &rawResp)
	return buildResponse(rawResp, err)
}

// Query returns a QueryBuilder to perform OData queries on the exposed view.
func (s *SQLView) Query() *QueryBuilder {
	return NewQuery(s.conn, "SQLViews('"+s.viewName+"')")
}

// CrossJoinBuilder helps in building crossjoin queries via QueryService_PostQuery.
type CrossJoinBuilder struct {
	conn    Executor
	path    string
	options []string
}

// NewCrossJoin initializes a query targeting QueryService_PostQuery for multiple tables.
func NewCrossJoin(conn Executor, tables ...string) *CrossJoinBuilder {
	return &CrossJoinBuilder{
		conn: conn,
		path: "$crossjoin(" + strings.Join(tables, ",") + ")",
	}
}

// Expand specifies the expand option for the cross join.
func (c *CrossJoinBuilder) Expand(expand string) *CrossJoinBuilder {
	c.options = append(c.options, "$expand="+expand)
	return c
}

// Filter specifies the filter option for the cross join.
func (c *CrossJoinBuilder) Filter(filter string) *CrossJoinBuilder {
	c.options = append(c.options, "$filter="+filter)
	return c
}

// Get executes the cross join query.
func (c *CrossJoinBuilder) Get() (*Response, error) {
	payload := map[string]string{
		"QueryPath": c.path,
	}
	if len(c.options) > 0 {
		payload["QueryOption"] = strings.Join(c.options, "&")
	}

	var rawResp map[string]any
	err := c.conn.Do("POST", "/QueryService_PostQuery", payload, &rawResp)
	return buildResponse(rawResp, err)
}

// buildResponse is a generic helper to construct a Response.
func buildResponse(rawResp map[string]any, err error) (*Response, error) {
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

	// Extract NextLink for pagination (Service Layer uses "odata.nextLink" or "@odata.nextLink")
	if nextLink, ok := rawResp["odata.nextLink"].(string); ok {
		resp.NextLink = nextLink
	} else if nextLink, ok = rawResp["@odata.nextLink"].(string); ok {
		resp.NextLink = nextLink
	}

	return resp, nil
}
