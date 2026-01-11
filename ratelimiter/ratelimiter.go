package ratelimiter

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Visitor struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	loclog := "[ratelimiter.NewRateLimiter]"
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     r,
		burst:    b,
	}

	go rl.cleanup()

	slog.Info(loclog, "info", "rate limiter initialized", "rate", r, "burst", b)
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &Visitor{
			Limiter:  limiter,
			LastSeen: time.Now(),
		}
		return limiter
	}

	v.LastSeen = time.Now()
	return v.Limiter
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.LastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
				slog.Info("[ratelimiter.cleanup]", "info", "visitor removed", "identifier", ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRequestIP(r)

		limiter := rl.getVisitor(ip)
		if !limiter.Allow() {
			slog.Warn("[ratelimiter]", "info", "rate limit exceeded", "identifier", ip)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}
