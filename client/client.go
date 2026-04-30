package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/ehitelrc/slsdk/errors"
)

// Client handles HTTP communications with the Service Layer.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient initializes a new HTTP client with a cookie jar for session management.
func NewClient(baseURL string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Jar: jar,
		},
	}
}

// Do executes an HTTP request, serializing/deserializing JSON automatically.
func (c *Client) Do(method, path string, reqBody any, resBody any) error {
	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to serialize request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error struct {
				Code    int `json:"code"`
				Message struct {
					Value string `json:"value"`
				} `json:"message"`
			} `json:"error"`
		}
		
		if parseErr := json.Unmarshal(bodyBytes, &errResp); parseErr == nil && errResp.Error.Message.Value != "" {
			return &errors.SAPError{
				Code:    errResp.Error.Code,
				Message: errResp.Error.Message.Value,
			}
		}
		
		// Fallback if parsing standard SAP error fails
		return &errors.SAPError{
			Code:    resp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	if resBody != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, resBody); err != nil {
			return fmt.Errorf("failed to deserialize response: %w", err)
		}
	}

	return nil
}
