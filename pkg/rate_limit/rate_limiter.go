package rate_limit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	b int
	// dur      time.Duration // clear key durarion
	interval time.Duration // rate limit durarion
	newLim   NewLim
	mu       *sync.RWMutex
	ticker   *time.Ticker
	mp       map[string]Lim
	exit     chan struct{}
}

func NewRL(secs int64, b int, newLim NewLim) (rl *RateLimiter, err error) {
	if secs < 1 || b < 1 || newLim == nil {
		return nil, fmt.Errorf("invlaid parameter for RateLimiter")
	}

	interval := time.Second * time.Duration(secs)

	rl = &RateLimiter{
		b:        b,
		interval: interval,
		newLim:   newLim,
		mu:       new(sync.RWMutex),
		ticker:   time.NewTicker(RATELIMITER_ClearEveryN * interval),
		mp:       make(map[string]Lim, 100),
		exit:     make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-rl.exit:
				// fmt.Println("!!! RateLimiter.ticker runtime closed")
				return
			case <-rl.ticker.C:
				// fmt.Println("~~~ RateLimiter.ticker", rfc3339now())
			}

			now := time.Now()
			rl.mu.Lock()
			defer rl.mu.Unlock()
			if len(rl.mp) <= 100 {
				return
			}

			for key, lim := range rl.mp {
				dur := now.Sub(lim.Last())
				if dur > RATELIMITER_ClearEveryN*rl.interval {
					// fmt.Println(rfc3339now(), "RateLimiter drop key:", key, dur)
					lim.Stop()
					delete(rl.mp, key)
				}
			}
		}
	}()

	return
}

func (rl *RateLimiter) GetLimiter(key string) (lim Lim) {
	var ok bool
	rl.mu.RLock()

	if lim, ok = rl.mp[key]; ok {
		rl.mu.RUnlock()
		return lim
	}

	rl.mu.RUnlock()
	rl.mu.Lock()

	if lim, ok = rl.mp[key]; !ok {
		lim = rl.newLim(rl.interval, rl.b)
		rl.mp[key] = lim
	}
	rl.mu.Unlock()

	return lim
}

func (rl *RateLimiter) Metrics() (time.Duration, int, int) {
	return rl.interval, rl.b, len(rl.mp)
}

func (rl *RateLimiter) Allow(key string) (ok bool) {
	return rl.GetLimiter(key).Allow(time.Now())
}

func (rl *RateLimiter) AllowWithContext(ctx context.Context, key string) (ok bool) {
	return rl.GetLimiter(key).AllowWithContext(ctx, time.Now())
}

func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.exit)

	rl.mu.Lock()
	for key, limiter := range rl.mp {
		limiter.Stop()
		delete(rl.mp, key)
	}
	rl.mu.Unlock()
}
