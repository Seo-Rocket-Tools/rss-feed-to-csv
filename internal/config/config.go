package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	// Server configuration
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	
	// RSS fetcher configuration
	RSSFetchTimeout time.Duration
	MaxRSSSize      int64 // TODO: Implement max RSS size check in fetcher
	UserAgent       string
	
	// Security configuration
	MaxURLLength    int
	RateLimitPerMin int
	
	// CSV export configuration  
	DefaultSanitize bool // TODO: Use as default when sanitize param not provided
	
	// Logging
	LogLevel string // TODO: Implement structured logging
}

// Load returns the application configuration from environment variables
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", ":8080"),
		ReadTimeout:     getDuration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getDuration("WRITE_TIMEOUT", 15*time.Second),
		ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		
		RSSFetchTimeout: getDuration("RSS_FETCH_TIMEOUT", 30*time.Second),
		MaxRSSSize:      getInt64("MAX_RSS_SIZE", 10*1024*1024), // 10MB
		UserAgent:       getEnv("USER_AGENT", "RSS-to-CSV-Exporter/1.0"),
		
		MaxURLLength:    getInt("MAX_URL_LENGTH", 2048),
		RateLimitPerMin: getInt("RATE_LIMIT_PER_MIN", 60),
		
		DefaultSanitize: getBool("DEFAULT_SANITIZE", false),
		LogLevel:        getEnv("LOG_LEVEL", "INFO"),
	}
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDuration gets a duration from environment variable
func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getInt gets an integer from environment variable
func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getInt64 gets an int64 from environment variable
func getInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBool gets a boolean from environment variable
func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}