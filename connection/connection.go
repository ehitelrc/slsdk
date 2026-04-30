package connection

import (
	"fmt"

	"github.com/ehitelrc/slsdk/client"
	"github.com/ehitelrc/slsdk/errors"
	"github.com/ehitelrc/slsdk/query"
)

// Config represents the credentials and settings required to connect to the Service Layer.
type Config struct {
	BaseURL  string
	Company  string
	UserName string
	Password string
}

// Connection manages the HTTP client and handles the session.
type Connection struct {
	Config Config
	client *client.Client
}

// NewConnection initializes a new Service Layer connection manager.
func NewConnection(cfg Config) *Connection {
	return &Connection{
		Config: cfg,
		client: client.NewClient(cfg.BaseURL),
	}
}

// loginRequest is the internal payload used for the POST /Login request.
type loginRequest struct {
	CompanyDB string `json:"CompanyDB"`
	UserName  string `json:"UserName"`
	Password  string `json:"Password"`
}

// Login authenticates against the Service Layer and stores the resulting session cookies.
func (c *Connection) Login() error {
	req := loginRequest{
		CompanyDB: c.Config.Company,
		UserName:  c.Config.UserName,
		Password:  c.Config.Password,
	}

	var res map[string]any
	err := c.client.Do("POST", "/Login", req, &res)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	return nil
}

// Do executes a generic HTTP request, with automatic retry for expired sessions.
func (c *Connection) Do(method, path string, reqBody any, resBody any) error {
	err := c.client.Do(method, path, reqBody, resBody)
	if err != nil && isSessionError(err) {
		// Attempt to auto-relogin
		if loginErr := c.Login(); loginErr == nil {
			// Retry the original request
			return c.client.Do(method, path, reqBody, resBody)
		}
	}
	return err
}

func isSessionError(err error) bool {
	if sapErr, ok := err.(*errors.SAPError); ok {
		// Code 301 is commonly used in SL for "Invalid Session"
		// Code 401 is standard HTTP Unauthorized
		if sapErr.Code == 301 || sapErr.Code == 401 {
			return true
		}
	}
	return false
}

// Items returns a query builder targeted at the Items (OITM) collection.
func (c *Connection) Items() *query.Builder {
	return query.NewBuilder("Items", c)
}
