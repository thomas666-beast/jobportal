package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// Limit defines a rate limit policy for one endpoint.
type Limit struct {
	// Requests per second sustained.
	Rate rate.Limit
	// Maximum burst size.
	Burst int
	// How long an idle bucket is kept before eviction.
	TTL time.Duration
}

// bucket tracks one client's limiter and last-seen time.
type bucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter is a per-IP token bucket store scoped to a single endpoint.
type IPRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	limit   Limit
}

func NewIPRateLimiter(limit Limit) *IPRateLimiter {
	if limit.TTL == 0 {
		limit.TTL = 10 * time.Minute
	}
	rl := &IPRateLimiter{
		buckets: make(map[string]*bucket),
		limit:   limit,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *IPRateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[ip]
	if !ok {
		b = &bucket{
			limiter: rate.NewLimiter(rl.limit.Rate, rl.limit.Burst),
		}
		rl.buckets[ip] = b
	}
	b.lastSeen = time.Now()
	return b.limiter
}

// cleanupLoop evicts idle buckets so memory doesn't grow unbounded.
func (rl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.limit.TTL / 2)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-rl.limit.TTL)
		rl.mu.Lock()
		for ip, b := range rl.buckets {
			if b.lastSeen.Before(cutoff) {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns an Echo middleware enforcing this limiter.
func (rl *IPRateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			if !rl.get(ip).Allow() {
				c.Response().Header().Set("Retry-After", "60")
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "too many requests, please slow down",
				})
			}

			return next(c)
		}
	}
}
