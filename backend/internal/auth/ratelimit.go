package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu           sync.Mutex
	ipRequests   map[string][]time.Time
	maxRequests  int
	perDuration  time.Duration
	cleanupEvery time.Duration
	lastCleanup  time.Time
}

// NewRateLimiter creates a new rate limiter
// maxRequests: maximum number of requests allowed
// perDuration: time window for maxRequests
// cleanupEvery: how often to clean old records
func NewRateLimiter(maxRequests int, perDuration, cleanupEvery time.Duration) *RateLimiter {
	return &RateLimiter{
		ipRequests:   make(map[string][]time.Time),
		maxRequests:  maxRequests,
		perDuration:  perDuration,
		cleanupEvery: cleanupEvery,
		lastCleanup:  time.Now(),
	}
}

// RateLimit creates middleware for rate limiting based on IP address
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Check if rate limit is exceeded
		if rl.isLimitExceeded(clientIP, now) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Clean up old records periodically
		rl.cleanupIfNeeded(now)

		// Process the request
		c.Next()
	}
}

// isLimitExceeded checks if the rate limit for an IP is exceeded
func (rl *RateLimiter) isLimitExceeded(ip string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Add current request time
	rl.ipRequests[ip] = append(rl.ipRequests[ip], now)

	// Remove requests that are outside the time window
	windowStart := now.Add(-rl.perDuration)
	var recentRequests []time.Time

	for _, t := range rl.ipRequests[ip] {
		if t.After(windowStart) {
			recentRequests = append(recentRequests, t)
		}
	}

	rl.ipRequests[ip] = recentRequests

	// Check if the number of recent requests exceeds the limit
	return len(recentRequests) > rl.maxRequests
}

// cleanupIfNeeded removes old records to prevent memory leaks
func (rl *RateLimiter) cleanupIfNeeded(now time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if cleanup is needed
	if now.Sub(rl.lastCleanup) < rl.cleanupEvery {
		return
	}

	// Perform cleanup
	windowStart := now.Add(-rl.perDuration)
	for ip, times := range rl.ipRequests {
		var recentRequests []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				recentRequests = append(recentRequests, t)
			}
		}

		if len(recentRequests) == 0 {
			delete(rl.ipRequests, ip)
		} else {
			rl.ipRequests[ip] = recentRequests
		}
	}

	rl.lastCleanup = now
}

// GetRequestCount returns the number of requests made by an IP in the current window
func (rl *RateLimiter) GetRequestCount(ip string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.perDuration)
	var count int

	for _, t := range rl.ipRequests[ip] {
		if t.After(windowStart) {
			count++
		}
	}

	return count
} 