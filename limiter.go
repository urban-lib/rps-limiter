package limiter

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Config struct {
	RPS   int           `mapstructure:"rps"`
	Burst int           `mapstructure:"burst"`
	TTL   time.Duration `mapstructure:"ttl"`
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	sync.RWMutex
	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
	ttl      time.Duration
}

func NewRPSLimiter(cfg *Config) *rateLimiter {
	return &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    rate.Limit(cfg.RPS),
		burst:    cfg.Burst,
		ttl:      cfg.TTL,
	}
}

func (r *rateLimiter) GetVisitor(ip string) *rate.Limiter {
	r.RLock()
	v, exists := r.visitors[ip]
	r.RUnlock()
	if !exists {
		limiter := rate.NewLimiter(r.limit, r.burst)
		r.Lock()
		r.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		r.Unlock()
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (r *rateLimiter) CleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		r.Lock()
		for ip, v := range r.visitors {
			if time.Since(v.lastSeen) > r.ttl {
				delete(r.visitors, ip)
			}
		}
		r.Unlock()
	}
}
