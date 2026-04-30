package objects

import "github.com/ehitelrc/slsdk/errors"

// Response is the unified response structure for operations.
type Response struct {
	Success bool
	Message string
	Data    any
	Error   *errors.SAPError
}
