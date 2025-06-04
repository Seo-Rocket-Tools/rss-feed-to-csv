package utils

import (
	"testing"
)

func TestHTMLSanitizer_StripHTML(t *testing.T) {
	s := NewHTMLSanitizer()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "Hello, World!",
			want:  "Hello, World!",
		},
		{
			name:  "simple HTML",
			input: "<p>Hello, <b>World</b>!</p>",
			want:  "Hello, World!",
		},
		{
			name:  "nested HTML",
			input: "<div><p>Nested <span>content</span></p></div>",
			want:  "Nested content",
		},
		{
			name:  "HTML entities",
			input: "Hello &amp; goodbye &lt;tag&gt;",
			want:  "Hello & goodbye",
		},
		{
			name:  "script tags",
			input: "<script>alert('xss')</script>Clean text",
			want:  "Clean text",
		},
		{
			name:  "multiple spaces",
			input: "<p>Multiple     spaces</p>",
			want:  "Multiple     spaces",
		},
		{
			name:  "line breaks",
			input: "<p>Line<br/>break</p>",
			want:  "Line\nbreak",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only HTML tags",
			input: "<div><span></span></div>",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.StripHTML(tt.input); got != tt.want {
				t.Errorf("StripHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}