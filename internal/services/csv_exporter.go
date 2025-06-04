package services

import (
	"context"
	"encoding/csv"
	"io"

	"rss-feed-to-csv/internal/models"
	"rss-feed-to-csv/internal/utils"
)

// CSVExporter handles exporting RSS data to CSV format
type CSVExporter struct {
	sanitizer *utils.HTMLSanitizer
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{
		sanitizer: utils.NewHTMLSanitizer(),
	}
}


// Export writes RSS items to CSV format
func (e *CSVExporter) Export(ctx context.Context, w io.Writer, rss *models.RSS, sanitizeHTML bool) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write headers
	headers := []string{"Title", "Link", "Description", "PubDate", "ImageURL", "Content"}
	if err := writer.Write(headers); err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	// Write data rows
	for _, item := range rss.Channel.Items {
		description := item.Description
		content := item.ContentEncoded
		
		// Apply HTML sanitization if requested
		if sanitizeHTML {
			description = e.sanitizer.StripHTML(description)
			content = e.sanitizer.StripHTML(content)
		}
		
		record := []string{
			item.Title,
			item.Link,
			description,
			item.PubDate,
			item.GetImageURL(),
			content,
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}