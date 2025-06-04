package validator

import (
	"strings"
	"testing"

	"rss-feed-to-csv/internal/errors"
)

func TestURLValidator_ValidateURL(t *testing.T) {
	v := NewURLValidator(2048)

	tests := []struct {
		name    string
		url     string
		wantErr bool
		errType error
	}{
		{
			name:    "valid http URL",
			url:     "http://example.com/rss",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com/feed.xml",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errType: errors.ErrEmptyURL,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com/rss",
			wantErr: true,
		},
		{
			name:    "missing host",
			url:     "https:///path",
			wantErr: true,
		},
		{
			name:    "URL too long",
			url:     "https://example.com/" + strings.Repeat("a", 2048),
			wantErr: true,
		},
		{
			name:    "malformed URL",
			url:     "not a url at all",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errType != nil && err != tt.errType {
				t.Errorf("ValidateURL() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestURLValidator_SanitizeInput(t *testing.T) {
	v := NewURLValidator(2048)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "https://example.com/rss",
			want:  "https://example.com/rss",
		},
		{
			name:  "HTML tags",
			input: "<script>alert('xss')</script>https://example.com",
			want:  "alert('xss')https://example.com",
		},
		{
			name:  "whitespace",
			input: "  https://example.com  ",
			want:  "https://example.com",
		},
		{
			name:  "multiple tags",
			input: "<b>bold</b> and <i>italic</i>",
			want:  "bold and italic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.SanitizeInput(tt.input); got != tt.want {
				t.Errorf("SanitizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}