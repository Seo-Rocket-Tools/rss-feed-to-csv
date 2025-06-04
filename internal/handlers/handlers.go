package handlers

import (
	"log"
	"net/http"

	"rss-feed-to-csv/internal/config"
	"rss-feed-to-csv/internal/services"
	"rss-feed-to-csv/internal/validator"
)

// Handler contains all HTTP handlers for the application
type Handler struct {
	rssFetcher  *services.RSSFetcher
	csvExporter *services.CSVExporter
	validator   *validator.URLValidator
}

// NewHandler creates a new handler with dependencies
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		rssFetcher:  services.NewRSSFetcher(cfg.RSSFetchTimeout, cfg.UserAgent),
		csvExporter: services.NewCSVExporter(),
		validator:   validator.NewURLValidator(cfg.MaxURLLength),
	}
}

// HandleIndex serves the main HTML page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

// HandleExport handles RSS to CSV export requests
func (h *Handler) HandleExport(w http.ResponseWriter, r *http.Request) {
	// Extract and validate query parameters
	rssURL := r.URL.Query().Get("url")
	rssURL = h.validator.SanitizeInput(rssURL)
	
	if err := h.validator.ValidateURL(rssURL); err != nil {
		log.Printf("[ERROR] Invalid URL - URL: %s, Error: %v, Client: %s", rssURL, err, r.RemoteAddr)
		http.Error(w, "Invalid URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	sanitize := r.URL.Query().Get("sanitize") == "true"
	
	log.Printf("[INFO] Fetching RSS feed - URL: %s, Client: %s, User-Agent: %s", 
		rssURL, r.RemoteAddr, r.Header.Get("User-Agent"))

	// Fetch RSS feed
	rss, err := h.rssFetcher.FetchRSS(r.Context(), rssURL)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch/parse RSS - URL: %s, Error: %v, Client: %s", 
			rssURL, err, r.RemoteAddr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Successfully parsed RSS feed - URL: %s, Items: %d, Sanitize: %v, Client: %s", 
		rssURL, len(rss.Channel.Items), sanitize, r.RemoteAddr)

	// Set response headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=feed.csv")

	// Export to CSV
	if err := h.csvExporter.Export(r.Context(), w, rss, sanitize); err != nil {
		log.Printf("[ERROR] Failed to export CSV - URL: %s, Error: %v, Client: %s", 
			rssURL, err, r.RemoteAddr)
		// Note: Headers already sent, can't return HTTP error
		return
	}

	log.Printf("[SUCCESS] CSV export completed - URL: %s, Items exported: %d, Client: %s", 
		rssURL, len(rss.Channel.Items), r.RemoteAddr)
}