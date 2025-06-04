package utils

import (
	"html"
	"regexp"
	"strings"
)

// HTMLSanitizer provides methods for sanitizing HTML content
type HTMLSanitizer struct {
	scriptRegex    *regexp.Regexp
	styleRegex     *regexp.Regexp
	brRegex        *regexp.Regexp
	blockRegex     *regexp.Regexp
	tagRegex       *regexp.Regexp
	multiLineRegex *regexp.Regexp
}

// NewHTMLSanitizer creates a new HTMLSanitizer with pre-compiled regex patterns
func NewHTMLSanitizer() *HTMLSanitizer {
	return &HTMLSanitizer{
		scriptRegex:    regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`),
		styleRegex:     regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`),
		brRegex:        regexp.MustCompile(`(?i)<br\s*/?>`),
		blockRegex:     regexp.MustCompile(`(?i)</?(p|div|h[1-6]|ul|ol|li)[^>]*>`),
		tagRegex:       regexp.MustCompile(`<[^>]+>`),
		multiLineRegex: regexp.MustCompile(`\n{3,}`),
	}
}

// StripHTML removes HTML tags and converts HTML entities to plain text
func (s *HTMLSanitizer) StripHTML(htmlStr string) string {
	if htmlStr == "" {
		return ""
	}

	// Decode HTML entities first
	htmlStr = html.UnescapeString(htmlStr)
	
	// Remove script and style elements
	htmlStr = s.scriptRegex.ReplaceAllString(htmlStr, "")
	htmlStr = s.styleRegex.ReplaceAllString(htmlStr, "")
	
	// Replace br tags with newlines
	htmlStr = s.brRegex.ReplaceAllString(htmlStr, "\n")
	
	// Replace block tags with newlines
	htmlStr = s.blockRegex.ReplaceAllString(htmlStr, "\n")
	
	// Remove all other HTML tags
	htmlStr = s.tagRegex.ReplaceAllString(htmlStr, "")
	
	// Clean up multiple newlines
	htmlStr = s.multiLineRegex.ReplaceAllString(htmlStr, "\n\n")
	
	// Trim whitespace
	htmlStr = strings.TrimSpace(htmlStr)
	
	return htmlStr
}