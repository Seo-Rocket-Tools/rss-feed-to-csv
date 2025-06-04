package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rss-feed-to-csv/internal/config"
	"rss-feed-to-csv/internal/handlers"
	"rss-feed-to-csv/internal/middleware"
)

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	
	// Load configuration
	cfg := config.Load()
	
	// Create handler with config
	handler := handlers.NewHandler(cfg)
	
	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitPerMin)
	
	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandleIndex)
	mux.HandleFunc("/export", rateLimiter.Limit(handler.HandleExport))
	
	// Create server with timeouts
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("[INFO] Server starting on http://localhost%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("[INFO] Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}
	
	log.Println("[INFO] Server stopped")
}