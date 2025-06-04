package errors

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrInvalidURL      = errors.New("invalid URL format")
	ErrEmptyURL        = errors.New("URL cannot be empty")
	ErrFetchTimeout    = errors.New("RSS feed fetch timeout")
	ErrInvalidRSSXML   = errors.New("invalid RSS XML format")
	ErrNoRSSItems      = errors.New("no items found in RSS feed")
	ErrCSVWriteFailed  = errors.New("failed to write CSV")
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// FetchError represents an error that occurred while fetching RSS
type FetchError struct {
	URL        string
	StatusCode int
	Err        error
}

func (e FetchError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("failed to fetch RSS from %s (status: %d): %v", e.URL, e.StatusCode, e.Err)
	}
	return fmt.Sprintf("failed to fetch RSS from %s: %v", e.URL, e.Err)
}

func (e FetchError) Unwrap() error {
	return e.Err
}