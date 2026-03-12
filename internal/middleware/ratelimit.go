package middleware

import (
	"net/http"
	"sync"
	"time"
)

// IPRequest tracks request counts for an IP
type IPRequest struct {
	count     int
	lastReset time.Time
}

var (
	mu       sync.Mutex
	vistors  = make(map[string]*IPRequest)
	limit    = 50
	duration = 10 * time.Second
)

// RateLimit middleware limits requests per IP
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		mu.Lock()
		v, exists := vistors[ip]
		if !exists {
			vistors[ip] = &IPRequest{count: 1, lastReset: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastReset) > duration {
			v.count = 1
			v.lastReset = time.Now()
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if v.count >= limit {
			mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		v.count++
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
