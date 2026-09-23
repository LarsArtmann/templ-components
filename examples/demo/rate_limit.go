package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Rate-limit knobs for the public demo (backlog #264): generous enough that
// no human demo flow ever trips them, tight enough that a scripted loop
// can't monopolize the instance. In-memory and per-instance by design — the
// demo holds no real data and runs single-instance.
const (
	demoRateLimit = 2.0 // requests per second, sustained, per client IP
	demoRateBurst = 20  // bucket depth: a full demo page's worth of calls
)

// clientBucket is one client's token bucket plus its last-seen time for
// janitor-free eviction.
type clientBucket struct {
	tokens float64
	last   time.Time
}

// ipLimiter is a per-IP token-bucket http.Handler middleware.
type ipLimiter struct {
	mu     sync.Mutex
	burst  float64
	perSec float64
	now    func() time.Time
	seen   map[string]*clientBucket
}

// newIPLimiter returns middleware allowing burst requests instantly, then
// refilling at perSec tokens per second per client IP.
func newIPLimiter(perSec float64, burst float64) func(http.Handler) http.Handler {
	limiter := &ipLimiter{
		burst:  burst,
		perSec: perSec,
		now:    time.Now,
		seen:   make(map[string]*clientBucket),
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.allow(clientIP(r)) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "demo rate limit exceeded", http.StatusTooManyRequests)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.seen[ip]
	if !ok {
		// Lazy eviction: when the map balloons past a sane bound, drop the
		// whole table. Clients re-bucket on their next request.
		if len(l.seen) > 4096 {
			l.seen = make(map[string]*clientBucket)
		}
		l.seen[ip] = &clientBucket{tokens: l.burst, last: now}

		return true
	}

	b.tokens += now.Sub(b.last).Seconds() * l.perSec
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--

	return true
}

// clientIP extracts the caller's IP from the request (RemoteAddr for the
// demo — Cloud Run terminates TLS at the load balancer, so X-Forwarded-For
// would be client-controlled spoofing bait for a limiter keyed on it).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
