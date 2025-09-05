package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	clients map[string]*rate.Limiter
	mu      sync.RWMutex
	rate    rate.Limit
	burst   int
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
	r := rate.Every(cfg.RateLimit.Window / time.Duration(cfg.RateLimit.Requests))
	return &RateLimiter{
		clients: make(map[string]*rate.Limiter),
		rate:    r,
		burst:   cfg.RateLimit.Requests,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.clients[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	go rl.cleanupRoutine()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getLimiter(ip)

			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, errors.ErrorResponse{
					Error: "Rate limit exceeded",
					Code:  "RATE_LIMIT_EXCEEDED",
				})
			}

			return next(c)
		}
	}
}

func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, limiter := range rl.clients {
			if limiter.Allow() {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}