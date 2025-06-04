package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"strings"
	"testing"

	"rss-feed-to-csv/internal/models"
)

func TestCSVExporter_Export(t *testing.T) {
	exporter := NewCSVExporter()
	ctx := context.Background()

	tests := []struct {
		name         string
		rss          *models.RSS
		sanitizeHTML bool
		wantHeaders  []string
		wantRows     int
		checkContent func(t *testing.T, records [][]string)
	}{
		{
			name: "basic export without sanitization",
			rss: &models.RSS{
				Channel: models.Channel{
					Items: []models.Item{
						{
							Title:          "Item 1",
							Link:           "https://example.com/1",
							Description:    "<p>Description 1</p>",
							PubDate:        "Mon, 02 Jan 2006",
							ContentEncoded: "<div>Content 1</div>",
						},
					},
				},
			},
			sanitizeHTML: false,
			wantHeaders:  []string{"Title", "Link", "Description", "PubDate", "ImageURL", "Content"},
			wantRows:     2, // headers + 1 item
			checkContent: func(t *testing.T, records [][]string) {
				if records[1][0] != "Item 1" {
					t.Errorf("Title = %q, want %q", records[1][0], "Item 1")
				}
				if records[1][2] != "<p>Description 1</p>" {
					t.Errorf("Description = %q, want %q", records[1][2], "<p>Description 1</p>")
				}
			},
		},
		{
			name: "export with HTML sanitization",
			rss: &models.RSS{
				Channel: models.Channel{
					Items: []models.Item{
						{
							Title:          "Item 1",
							Link:           "https://example.com/1",
							Description:    "<p>Description <b>bold</b></p>",
							PubDate:        "Mon, 02 Jan 2006",
							ContentEncoded: "<script>alert('xss')</script>Clean content",
						},
					},
				},
			},
			sanitizeHTML: true,
			wantHeaders:  []string{"Title", "Link", "Description", "PubDate", "ImageURL", "Content"},
			wantRows:     2,
			checkContent: func(t *testing.T, records [][]string) {
				if !strings.Contains(records[1][2], "Description bold") {
					t.Errorf("Sanitized description = %q, should contain 'Description bold'", records[1][2])
				}
				if strings.Contains(records[1][5], "<script>") {
					t.Errorf("Content contains script tag after sanitization: %q", records[1][5])
				}
				if !strings.Contains(records[1][5], "Clean content") {
					t.Errorf("Content = %q, should contain 'Clean content'", records[1][5])
				}
			},
		},
		{
			name: "export with image URLs",
			rss: &models.RSS{
				Channel: models.Channel{
					Items: []models.Item{
						{
							Title: "Item with image",
							MediaContent: []models.MediaContent{
								{
									URL:    "https://example.com/image.jpg",
									Type:   "image/jpeg",
									Medium: "image",
								},
							},
						},
						{
							Title: "Item with media",
							MediaContent: []models.MediaContent{
								{
									URL:  "https://example.com/media.png",
									Type: "image/png",
								},
							},
						},
					},
				},
			},
			sanitizeHTML: false,
			wantRows:     3, // headers + 2 items
			checkContent: func(t *testing.T, records [][]string) {
				if records[1][4] != "https://example.com/image.jpg" {
					t.Errorf("ImageURL[0] = %q, want %q", records[1][4], "https://example.com/image.jpg")
				}
				if records[2][4] != "https://example.com/media.png" {
					t.Errorf("ImageURL[1] = %q, want %q", records[2][4], "https://example.com/media.png")
				}
			},
		},
		{
			name: "empty RSS feed",
			rss: &models.RSS{
				Channel: models.Channel{
					Items: []models.Item{},
				},
			},
			sanitizeHTML: false,
			wantRows:     1, // headers only
		},
		{
			name: "multiple items",
			rss: &models.RSS{
				Channel: models.Channel{
					Items: []models.Item{
						{Title: "Item 1", Link: "https://example.com/1"},
						{Title: "Item 2", Link: "https://example.com/2"},
						{Title: "Item 3", Link: "https://example.com/3"},
					},
				},
			},
			sanitizeHTML: false,
			wantRows:     4, // headers + 3 items
			checkContent: func(t *testing.T, records [][]string) {
				for i := 1; i <= 3; i++ {
					if records[i][0] != strings.TrimSpace("Item "+string(rune('0'+i))) {
						t.Errorf("Title[%d] = %q, want %q", i-1, records[i][0], "Item "+string(rune('0'+i)))
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := exporter.Export(ctx, &buf, tt.rss, tt.sanitizeHTML)
			if err != nil {
				t.Fatalf("Export() error = %v", err)
			}

			// Parse the CSV output
			reader := csv.NewReader(&buf)
			records, err := reader.ReadAll()
			if err != nil {
				t.Fatalf("Failed to read CSV: %v", err)
			}

			if len(records) != tt.wantRows {
				t.Errorf("Number of rows = %d, want %d", len(records), tt.wantRows)
			}

			if len(records) > 0 && tt.wantHeaders != nil {
				headers := records[0]
				for i, want := range tt.wantHeaders {
					if i >= len(headers) || headers[i] != want {
						t.Errorf("Header[%d] = %q, want %q", i, headers[i], want)
					}
				}
			}

			if tt.checkContent != nil {
				tt.checkContent(t, records)
			}
		})
	}
}

func TestCSVExporter_Export_WriterError(t *testing.T) {
	exporter := NewCSVExporter()
	ctx := context.Background()

	// Create a writer that fails immediately
	failWriter := &failingWriter{failAfter: 0}

	rss := &models.RSS{
		Channel: models.Channel{
			Items: []models.Item{
				{Title: "Item 1"},
			},
		},
	}

	err := exporter.Export(ctx, failWriter, rss, false)
	if err == nil {
		t.Error("Export() should return error when writer fails")
	}
}

// Define a test error
var testWriteError = errors.New("test write error")

// failingWriter is a test writer that fails after n writes
type failingWriter struct {
	writes    int
	failAfter int
}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	if w.writes >= w.failAfter {
		return 0, testWriteError
	}
	w.writes++
	return len(p), nil
}