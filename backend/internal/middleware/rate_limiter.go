package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
)

type bucket struct {
	count   int
	resetAt time.Time
}

// RateLimiter is a per-IP token bucket.
type RateLimiter struct {
	mu     sync.Mutex
	limits map[string]*bucket
	reqs   int
	window time.Duration
}

// NewRateLimiter creates a limiter.
func NewRateLimiter(reqs int, window time.Duration) *RateLimiter {
	return &RateLimiter{limits: make(map[string]*bucket), reqs: reqs, window: window}
}

// Allow records one request for ip and reports whether it is within the limit.
func (r *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	b, ok := r.limits[ip]
	if !ok || now.After(b.resetAt) {
		b = &bucket{count: 0, resetAt: now.Add(r.window)}
		r.limits[ip] = b
	}
	b.count++
	return b.count <= r.reqs
}

// Snapshot returns per-IP request counts for observability.
func (r *RateLimiter) Snapshot() map[string]int {
	out := make(map[string]int, len(r.limits))
	for ip, b := range r.limits {
		out[ip] = b.count
	}
	return out
}

// Reset clears the counter for one ip.
func (r *RateLimiter) Reset(ip string) {
	delete(r.limits, ip)
}

// Limit returns a middleware enforcing the rate limit per client IP.
func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !r.Allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				dto.Fail(constants.CodeRateLimited, constants.MsgRateLimited))
			return
		}
		c.Next()
	}
}
