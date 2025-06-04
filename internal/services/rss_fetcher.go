package services

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"rss-feed-to-csv/internal/errors"
	"rss-feed-to-csv/internal/models"
)

// RSSFetcher handles fetching and parsing RSS feeds
type RSSFetcher struct {
	client    *http.Client
	userAgent string
}

// NewRSSFetcher creates a new RSS fetcher with configured HTTP client
func NewRSSFetcher(timeout time.Duration, userAgent string) *RSSFetcher {
	return &RSSFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: userAgent,
	}
}

// FetchRSS fetches and parses an RSS feed from the given URL
func (f *RSSFetcher) FetchRSS(ctx context.Context, url string) (*models.RSS, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers to improve compatibility with RSS servers
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml, */*")
	
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.FetchError{
			URL:        url,
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("unexpected status: %s", resp.Status),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS feed: %w", err)
	}

	var rss models.RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	// Validate RSS has content
	if len(rss.Channel.Items) == 0 {
		return nil, errors.ErrNoRSSItems
	}

	return &rss, nil
}