package errors

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "url",
		Message: "invalid format",
	}

	expected := "validation error on field 'url': invalid format"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestFetchError(t *testing.T) {
	tests := []struct {
		name     string
		err      FetchError
		expected string
	}{
		{
			name: "with status code",
			err: FetchError{
				URL:        "https://example.com/feed",
				StatusCode: 404,
				Err:        errors.New("not found"),
			},
			expected: "failed to fetch RSS from https://example.com/feed (status: 404): not found",
		},
		{
			name: "without status code",
			err: FetchError{
				URL: "https://example.com/feed",
				Err: errors.New("connection timeout"),
			},
			expected: "failed to fetch RSS from https://example.com/feed: connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFetchError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	err := FetchError{
		URL: "https://example.com",
		Err: innerErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}

func TestCommonErrors(t *testing.T) {
	// Test that common errors are defined
	commonErrors := []struct {
		err  error
		desc string
	}{
		{ErrInvalidURL, "invalid URL format"},
		{ErrEmptyURL, "URL cannot be empty"},
		{ErrFetchTimeout, "RSS feed fetch timeout"},
		{ErrInvalidRSSXML, "invalid RSS XML format"},
		{ErrNoRSSItems, "no items found in RSS feed"},
		{ErrCSVWriteFailed, "failed to write CSV"},
	}

	for _, ce := range commonErrors {
		if ce.err == nil {
			t.Errorf("Common error %q is nil", ce.desc)
		}
		if ce.err.Error() != ce.desc {
			t.Errorf("Common error = %q, want %q", ce.err.Error(), ce.desc)
		}
	}
}