package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalEnv := make(map[string]string)
	envVars := []string{
		"PORT", "READ_TIMEOUT", "WRITE_TIMEOUT", "SHUTDOWN_TIMEOUT",
		"RSS_FETCH_TIMEOUT", "MAX_RSS_SIZE", "USER_AGENT",
		"MAX_URL_LENGTH", "RATE_LIMIT_PER_MIN", "DEFAULT_SANITIZE", "LOG_LEVEL",
	}
	
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	
	// Restore env vars after test
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("default values", func(t *testing.T) {
		cfg := Load()
		
		if cfg.Port != ":8080" {
			t.Errorf("Port = %s, want :8080", cfg.Port)
		}
		if cfg.ReadTimeout != 15*time.Second {
			t.Errorf("ReadTimeout = %v, want 15s", cfg.ReadTimeout)
		}
		if cfg.RSSFetchTimeout != 30*time.Second {
			t.Errorf("RSSFetchTimeout = %v, want 30s", cfg.RSSFetchTimeout)
		}
		if cfg.MaxRSSSize != 10*1024*1024 {
			t.Errorf("MaxRSSSize = %d, want %d", cfg.MaxRSSSize, 10*1024*1024)
		}
		if cfg.UserAgent != "RSS-to-CSV-Exporter/1.0" {
			t.Errorf("UserAgent = %s, want RSS-to-CSV-Exporter/1.0", cfg.UserAgent)
		}
		if cfg.MaxURLLength != 2048 {
			t.Errorf("MaxURLLength = %d, want 2048", cfg.MaxURLLength)
		}
		if cfg.DefaultSanitize != false {
			t.Errorf("DefaultSanitize = %v, want false", cfg.DefaultSanitize)
		}
	})

	t.Run("custom values from env", func(t *testing.T) {
		os.Setenv("PORT", ":9090")
		os.Setenv("READ_TIMEOUT", "20s")
		os.Setenv("MAX_RSS_SIZE", "5242880")
		os.Setenv("DEFAULT_SANITIZE", "true")
		os.Setenv("LOG_LEVEL", "DEBUG")
		
		cfg := Load()
		
		if cfg.Port != ":9090" {
			t.Errorf("Port = %s, want :9090", cfg.Port)
		}
		if cfg.ReadTimeout != 20*time.Second {
			t.Errorf("ReadTimeout = %v, want 20s", cfg.ReadTimeout)
		}
		if cfg.MaxRSSSize != 5242880 {
			t.Errorf("MaxRSSSize = %d, want 5242880", cfg.MaxRSSSize)
		}
		if cfg.DefaultSanitize != true {
			t.Errorf("DefaultSanitize = %v, want true", cfg.DefaultSanitize)
		}
		if cfg.LogLevel != "DEBUG" {
			t.Errorf("LogLevel = %s, want DEBUG", cfg.LogLevel)
		}
	})

	t.Run("invalid env values use defaults", func(t *testing.T) {
		os.Setenv("READ_TIMEOUT", "invalid")
		os.Setenv("MAX_RSS_SIZE", "not-a-number")
		os.Setenv("DEFAULT_SANITIZE", "not-a-bool")
		
		cfg := Load()
		
		// Should fall back to defaults
		if cfg.ReadTimeout != 15*time.Second {
			t.Errorf("ReadTimeout = %v, want 15s (default)", cfg.ReadTimeout)
		}
		if cfg.MaxRSSSize != 10*1024*1024 {
			t.Errorf("MaxRSSSize = %d, want %d (default)", cfg.MaxRSSSize, 10*1024*1024)
		}
		if cfg.DefaultSanitize != false {
			t.Errorf("DefaultSanitize = %v, want false (default)", cfg.DefaultSanitize)
		}
	})
}