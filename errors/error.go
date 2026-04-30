package errors

import "fmt"

// SAPError represents an error returned by the SAP Business One Service Layer.
type SAPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *SAPError) Error() string {
	return fmt.Sprintf("SAP Error %d: %s", e.Code, e.Message)
}
