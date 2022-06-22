package rate_limit

import (
	"context"
	"fmt"
	"time"
)

// rate limiter using buffered channel
type LimiterV2 struct {
	ch     chan bool
	ticker *time.Ticker
	last   time.Time
	exit   chan struct{}
}

func NewLimiterV2(interval time.Duration, b int) (limiter *LimiterV2, err error) {
	if interval < time.Second || b < 1 {
		return limiter, fmt.Errorf("invlaid parameter for NewLimiter")
	}

	limiter = &LimiterV2{
		ch:     make(chan bool, b),
		ticker: time.NewTicker(interval),
		exit:   make(chan struct{}),
	}

	for i := 0; i < b; i++ {
		limiter.ch <- true
	}

	go func() {
		// _, ok := <-limiter.ticker.C; !ok // !! time.Ticker.Stop doesn't close the channel
		for {
			select {
			case <-limiter.exit:
				// fmt.Println("!!! Limiter.ticker runtime closed", rfc3339now())
				// close(limiter.ch)
				return
			case <-limiter.ticker.C:
				// fmt.Println("~~~ Limiter.ticker", rfc3339now())
			}
		loop:
			for i := 0; i < b; i++ {
				select {
				case <-limiter.exit:
					return
				case limiter.ch <- true:
					// fmt.Println(rfc3339now(), "add token to bucket")
				default:
					// fmt.Println(rfc3339now(), "bucket is full")
					break loop
				}
			}
		}
	}()

	return
}

func (limiter *LimiterV2) New(interval time.Duration, b int) (*LimiterV2, error) {
	return NewLimiterV2(interval, b)
}

func (limiter *LimiterV2) Allow(now time.Time) (ok bool) {
	if limiter.last.After(now) { // now.IsZero()
		now = time.Now()
	}
	limiter.last = now

	select {
	case <-limiter.exit:
		return false
	case _, ok = <-limiter.ch: // ?? will limiter.ch be closed
		return ok
	default:
		return false // return immediately
	}
}

func (limiter *LimiterV2) AllowWithContext(ctx context.Context, now time.Time) (ok bool) {
	if limiter.last.After(now) { // now.IsZero()
		now = time.Now()
	}
	limiter.last = now

	select {
	case <-limiter.exit:
		return false
	case <-ctx.Done():
		return false
	case _, ok = <-limiter.ch: // ?? will limiter.ch be closed
		return ok
	}
}

func (limiter *LimiterV2) Last() time.Time {
	return limiter.last
}

func (limiter *LimiterV2) Stop() {
	limiter.ticker.Stop()
	close(limiter.exit)
}
