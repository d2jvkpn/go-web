package rate_limit

import (
	"context"
	"time"
)

const (
	RATELIMITER_ClearEveryN = 5
)

func rfc3339now() string {
	return time.Now().Format(time.RFC3339)
}

type Lim interface {
	Allow(time.Time) bool
	AllowWithContext(context.Context, time.Time) bool
	Stop()
	Last() time.Time
}

type NewLim = func(time.Duration, int) Lim

func NewLimiter(interval time.Duration, b int) (lim Lim, err error) {
	return NewLimiterV3(interval, b, false)
}

func NewRateLimiter(secs int64, b int) (rl *RateLimiter, err error) {
	newLimi := func(dur time.Duration, b int) Lim {
		limiter, _ := NewLimiterV2(dur, b)
		return limiter
	}

	return NewRL(secs, b, newLimi)
}
