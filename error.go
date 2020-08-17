package monobank

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// ReqError request error
type ReqError struct {
	Method string
	URL    *url.URL
	Err    error // wrapped error
}

func NewReqError(req *http.Request, err error) error {
	if req == nil {
		return errors.New("empty request")
	}

	return &ReqError{
		Method: req.Method,
		URL:    req.URL,
		Err:    err,
	}
}

func (e *ReqError) Error() string {
	return fmt.Sprintf("request %s %s: %s", e.Method, e.URL, e.Err.Error())
}

// Cause is causer interface(https://github.com/pkg/errors)
func (e *ReqError) Cause() error {
	return e.Err
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *ReqError) Unwrap() error {
	return e.Err
}

// APIError monobank API errors
type APIError struct {
	ResponseStatusCode  int // HTTP status code
	ExpectedStatusCodes []int
	Err                 error // wrapped error, a body as for now
}

func (e *APIError) Error() string {
	return fmt.Sprintf("unexpected status code %d(want %v): %s",
		e.ResponseStatusCode, e.ExpectedStatusCodes, e.Err)
}

func (e *APIError) Cause() error {
	return e.Err
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *APIError) Unwrap() error {
	return e.Err
}
