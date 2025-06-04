package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	cleanup  time.Duration // cleanup interval for old visitors
}

// visitor tracks the rate limit state for each client
type visitor struct {
	tokens    int
	lastVisit time.Time
	mu        sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(ratePerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     ratePerMinute,
		cleanup:  5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// Limit returns a middleware that enforces rate limiting
func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		// Get or create visitor
		v := rl.getVisitor(ip)

		// Check if visitor has tokens
		if !v.allowRequest() {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

// getVisitor retrieves or creates a visitor for the given IP
func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.RLock()
	v, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		v, exists = rl.visitors[ip]
		if !exists {
			v = &visitor{
				tokens:    rl.rate,
				lastVisit: time.Now(),
			}
			rl.visitors[ip] = v
		}
		rl.mu.Unlock()
	}

	return v
}

// allowRequest checks if a request is allowed and updates tokens
func (v *visitor) allowRequest() bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.lastVisit)

	// Refill tokens based on elapsed time
	// We get 1 token per second (60 per minute)
	tokensToAdd := int(elapsed.Seconds())
	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > 60 { // Cap at rate per minute
			v.tokens = 60
		}
		v.lastVisit = now
	}

	// Check if we have tokens
	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// cleanupVisitors removes old visitor entries
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if now.Sub(v.lastVisit) > rl.cleanup {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Take the first IP if there are multiple
		if idx := indexOf(ip, ','); idx != -1 {
			ip = ip[:idx]
		}
		return ip
	}

	// Check X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	if idx := lastIndexOf(r.RemoteAddr, ':'); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// indexOf returns the index of the first occurrence of sep in s
func indexOf(s string, sep rune) int {
	for i, r := range s {
		if r == sep {
			return i
		}
	}
	return -1
}

// lastIndexOf returns the index of the last occurrence of sep in s
func lastIndexOf(s string, sep rune) int {
	lastIdx := -1
	for i, r := range s {
		if r == sep {
			lastIdx = i
		}
	}
	return lastIdx
}