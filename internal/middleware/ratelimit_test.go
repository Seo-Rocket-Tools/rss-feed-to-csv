package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRateLimiter_Limit(t *testing.T) {
	// Create a rate limiter with 10 requests per minute
	rl := NewRateLimiter(10)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap the handler with rate limiting
	handler := rl.Limit(testHandler)

	// Test allowing initial requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// Test rate limit exceeded
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	// Create a rate limiter
	rl := NewRateLimiter(60) // 60 per minute = 1 per second

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(testHandler)

	// Use up one token
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:1234"
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("First request failed: %d", w.Code)
	}

	// Wait for token refill (2 seconds to be safe)
	time.Sleep(2 * time.Second)

	// Should be able to make another request
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request after refill failed: %d", w.Code)
	}
}

func TestRateLimiter_MultipleClients(t *testing.T) {
	rl := NewRateLimiter(5)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(testHandler)

	clients := []string{
		"192.168.1.1:1234",
		"192.168.1.2:1234",
		"192.168.1.3:1234",
	}

	// Each client should be able to make 5 requests
	for _, client := range clients {
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = client
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Client %s request %d failed: %d", client, i+1, w.Code)
			}
		}
	}

	// Each client should be rate limited on the 6th request
	for _, client := range clients {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = client
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Client %s should be rate limited, got: %d", client, w.Code)
		}
	}
}

func TestRateLimiter_Concurrent(t *testing.T) {
	rl := NewRateLimiter(100)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(testHandler)

	var wg sync.WaitGroup
	successCount := 0
	rateLimitCount := 0
	var mu sync.Mutex

	// Run 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			w := httptest.NewRecorder()

			handler(w, req)

			mu.Lock()
			defer mu.Unlock()

			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				rateLimitCount++
			}
		}()
	}

	wg.Wait()

	// Should have approximately 100 successful requests
	if successCount > 100 {
		t.Errorf("Too many successful requests: %d", successCount)
	}

	if rateLimitCount < 50 {
		t.Errorf("Too few rate limited requests: %d", rateLimitCount)
	}

	if successCount+rateLimitCount != 200 {
		t.Errorf("Total requests don't match: success=%d, limited=%d", successCount, rateLimitCount)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name        string
		setupReq    func(*http.Request)
		expectedIP  string
	}{
		{
			name: "RemoteAddr only",
			setupReq: func(r *http.Request) {
				r.RemoteAddr = "192.168.1.1:1234"
			},
			expectedIP: "192.168.1.1",
		},
		{
			name: "X-Forwarded-For single IP",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Forwarded-For", "10.0.0.1")
				r.RemoteAddr = "192.168.1.1:1234"
			},
			expectedIP: "10.0.0.1",
		},
		{
			name: "X-Forwarded-For multiple IPs",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2, 10.0.0.3")
				r.RemoteAddr = "192.168.1.1:1234"
			},
			expectedIP: "10.0.0.1",
		},
		{
			name: "X-Real-IP",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Real-IP", "10.0.0.2")
				r.RemoteAddr = "192.168.1.1:1234"
			},
			expectedIP: "10.0.0.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupReq(req)

			ip := getClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}