package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	visitors map[string]*rateLimiter
	mu       sync.RWMutex
}

func NewRateLimiterStore() *RateLimiterStore {
	store := &RateLimiterStore{
		visitors: make(map[string]*rateLimiter),
	}

	go store.cleanupVisitors()
	return store
}

func (store *RateLimiterStore) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		store.mu.Lock()
		for ip, v := range store.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(store.visitors, ip)
			}
		}
		store.mu.Unlock()
	}
}

func (store *RateLimiterStore) getRateLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	v, exists := store.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(r, b)
		store.visitors[ip] = &rateLimiter{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func RateLimitMiddleware(requestsPerSecond float64, burstSize int) echo.MiddlewareFunc {
	store := NewRateLimiterStore()
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := store.getRateLimiter(ip, rate.Limit(requestsPerSecond), burstSize)
			
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}
			
			return next(c)
		}
	}
}