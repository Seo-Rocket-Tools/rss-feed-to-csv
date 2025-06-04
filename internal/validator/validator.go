package validator

import (
	"net/url"
	"regexp"
	"strings"

	"rss-feed-to-csv/internal/errors"
)

// URLValidator validates URLs for RSS feeds
type URLValidator struct {
	allowedSchemes map[string]bool
	maxURLLength   int
}

// NewURLValidator creates a new URL validator
func NewURLValidator(maxURLLength int) *URLValidator {
	return &URLValidator{
		allowedSchemes: map[string]bool{
			"http":  true,
			"https": true,
		},
		maxURLLength: maxURLLength,
	}
}

// ValidateURL validates a URL for RSS fetching
func (v *URLValidator) ValidateURL(rawURL string) error {
	if rawURL == "" {
		return errors.ErrEmptyURL
	}

	if len(rawURL) > v.maxURLLength {
		return &errors.ValidationError{
			Field:   "url",
			Message: "URL is too long",
		}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &errors.ValidationError{
			Field:   "url",
			Message: "invalid URL format",
		}
	}

	if !v.allowedSchemes[parsedURL.Scheme] {
		return &errors.ValidationError{
			Field:   "url", 
			Message: "URL scheme must be http or https",
		}
	}

	if parsedURL.Host == "" {
		return &errors.ValidationError{
			Field:   "url",
			Message: "URL must have a host",
		}
	}

	return nil
}

// SanitizeInput sanitizes user input to prevent XSS
func (v *URLValidator) SanitizeInput(input string) string {
	// Basic sanitization - remove any HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	sanitized := re.ReplaceAllString(input, "")
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

