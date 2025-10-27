package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

const FiveRequestsPerMinute = 5
const TenMinutesThreshold = 10

// Token Bucket Algorithm
// Each IP address has a bucket that holds tokens. Each request costs 1 token, and tokens refill continuously over time.
type bucket struct {
	tokens         float64
	lastRefillTime time.Time
	mu             sync.Mutex
}

type IpThrottler struct {
	buckets    map[string]*bucket
	mu         sync.RWMutex
	maxTokens  float64
	refillRate float64
}

// NewIpThrottler initializes and returns a new IpThrottler with the specified maximum requests per second.
func NewIpThrottler(maxRequests int) *IpThrottler {
	throttler := &IpThrottler{
		buckets:    make(map[string]*bucket),
		maxTokens:  float64(maxRequests),
		refillRate: float64(maxRequests),
	}

	go throttler.cleanup()
	return throttler
}

// cleanup removes inactive buckets from the IpThrottler's map if they have not been accessed within the last 10 minutes.
func (t *IpThrottler) cleanup() {
	ticker := time.NewTicker(FiveRequestsPerMinute * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		t.mu.Lock()
		now := time.Now()

		for ip, bu := range t.buckets {
			bu.mu.Lock()

			if now.Sub(bu.lastRefillTime) > TenMinutesThreshold*time.Minute {
				delete(t.buckets, ip)
			}
			bu.mu.Unlock()
		}
		t.mu.Unlock()
	}
}

// getBucket retrieves the bucket associated with the given IP address or creates a new one if it does not exist.
func (t *IpThrottler) getBucket(ip string) *bucket {
	t.mu.RLock()
	b, exists := t.buckets[ip]
	t.mu.RUnlock()

	if exists {
		return b
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if b, exists := t.buckets[ip]; exists {
		return b
	}

	b = &bucket{
		tokens:         t.maxTokens,
		lastRefillTime: time.Now(),
	}
	t.buckets[ip] = b
	return b
}

// allowRequests determines if a request from the given IP address should be allowed based on token-bucket rate limiting.
func (t *IpThrottler) allowRequests(ip string) bool {
	b := t.getBucket(ip)
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefillTime).Seconds()

	b.tokens += elapsed * t.refillRate

	if b.tokens > t.maxTokens {
		b.tokens = t.maxTokens
	}

	b.lastRefillTime = now

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}

	return false
}

// Throttler is a middleware function that limits the number of requests per second from a single client IP address.
func Throttler(maxRequestsPerSecond int) gin.HandlerFunc {
	throttler := NewIpThrottler(maxRequestsPerSecond)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !throttler.allowRequests(ip) {
			utils.RespondWithError(c, http.StatusTooManyRequests, "Too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
