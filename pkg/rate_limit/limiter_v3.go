package rate_limit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// rate limiter using a time array
type LimiterV3 struct {
	interval time.Duration
	vec      []time.Time
	strong   bool // count event get bucket failed
	p        int  // index of last updated
	mu       sync.RWMutex
	exit     chan struct{}
}

func NewLimiterV3(interval time.Duration, b int, strong bool) (limiter *LimiterV3, err error) {
	if interval < time.Second || b <= 0 {
		return nil, fmt.Errorf("invalid parameter for NewLimiter")
	}

	limiter = &LimiterV3{
		interval: interval,
		vec:      make([]time.Time, b),
		strong:   strong,
		exit:     make(chan struct{}),
	}

	return
}

// (intf) New(time.Duration, int, bool) (intf, error), as golang 1.17 doesn't support generics
func (limiter *LimiterV3) New(interval time.Duration, b int, strong bool) (*LimiterV3, error) {
	return NewLimiterV3(interval, b, strong)
}

func (limiter *LimiterV3) next(now time.Time) (p, next int) {
	p = limiter.p
	switch {
	case limiter.p == 0 && limiter.vec[0].IsZero():
		next = 0
	case limiter.p < len(limiter.vec)-1:
		next = limiter.p + 1
	default:
		next = 0
	}

	return p, next
}

func (limiter *LimiterV3) update(p, next int, now time.Time) (ok bool) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	if limiter.p == p { // not changed by other concurrent routines
		limiter.p, limiter.vec[next] = next, now
		return true
	}

	// check again as limiter.p changed by other concurrent routines already
	_, next2 := limiter.next(now)
	if ok = now.Sub(limiter.vec[next2]) > limiter.interval; !ok {
		return false
	}

	limiter.p, limiter.vec[next2] = next2, now
	return true
}

func (limiter *LimiterV3) allow(now time.Time) (ok bool) {
	if limiter.vec[limiter.p].After(now) {
		now = time.Now()
	}

	limiter.mu.RLock()
	p, next := limiter.next(now)
	limiter.mu.RUnlock()

	ok = now.Sub(limiter.vec[next]) > limiter.interval

	if limiter.strong || ok {
		// limiter.p, limiter.vec[next] = next, now
		ok = limiter.update(p, next, now)
	}

	return ok
}

func (limiter *LimiterV3) Allow(now time.Time) (ok bool) {
	select {
	case <-limiter.exit:
		return false
	default:
	}

	return limiter.allow(now)
}

func (limiter *LimiterV3) AllowWithContext(ctx context.Context, now time.Time) (ok bool) {
	select {
	case <-limiter.exit:
		return false
	case <-ctx.Done(): // allow context like a timeout
		return false
	default:
	}

	return limiter.allow(now)
}

func (limiter *LimiterV3) Last() time.Time {
	return limiter.vec[limiter.p]
}

func (limiter *LimiterV3) Stop() {
	close(limiter.exit)
}
