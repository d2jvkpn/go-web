package rate_limit

import (
	"context"
	"fmt"
	"time"
)

// rate limiter using a time array
type LimiterV1 struct {
	interval time.Duration
	vec      []time.Time
	strong   bool // count event get bucket failed
	p        int
	ch       chan struct{}
	exit     chan struct{}
}

func NewLimiterV1(interval time.Duration, b int, strong bool) (limiter *LimiterV1, err error) {
	if interval < time.Second || b <= 0 {
		return nil, fmt.Errorf("invalid parameter for NewLimiter")
	}

	limiter = &LimiterV1{
		interval: interval,
		vec:      make([]time.Time, b),
		strong:   strong,
		ch:       make(chan struct{}, 1),
		exit:     make(chan struct{}),
	}

	return
}

// (intf) New(time.Duration, int, bool) (intf, error), as golang 1.17 doesn't support generics
func (limiter *LimiterV1) New(interval time.Duration, b int, strong bool) (*LimiterV1, error) {
	return NewLimiterV1(interval, b, strong)
}

//func (limiter *LimiterV1) pNext(now time.Time) (oldest time.Time) {
//	switch {
//	case limiter.p == 0 && limiter.vec[0].IsZero():
//	case limiter.p < len(limiter.vec)-1:
//		limiter.p++
//	default:
//		limiter.p = 0
//	}

//	oldest = limiter.vec[limiter.p] // extract oldest value

//	return oldest
//}

//func (limiter *LimiterV1) pBack() {
//	switch {
//	case limiter.p == 0 && limiter.vec[0].IsZero():
//	case limiter.p > 0:
//		limiter.p--
//	default:
//		limiter.p = len(limiter.vec) - 1
//	}
//}

func (limiter *LimiterV1) next(now time.Time) (next int) {
	switch {
	case limiter.p == 0 && limiter.vec[0].IsZero():
		next = 0
	case limiter.p < len(limiter.vec)-1:
		next = limiter.p + 1
	default:
		next = 0
	}

	return next
}

func (limiter *LimiterV1) allow(now time.Time) (ok bool) {
	if limiter.vec[limiter.p].After(now) {
		now = time.Now()
	}

	//	ok = now.Sub(limiter.pNext(now)) > limiter.interval

	//	if limiter.strong || ok {
	//		limiter.vec[limiter.p] = now
	//	} else {
	//		limiter.pBack()
	//	}

	next := limiter.next(now)
	ok = now.Sub(limiter.vec[next]) > limiter.interval
	if limiter.strong || ok {
		limiter.p, limiter.vec[next] = next, now
	}

	return ok
}

func (limiter *LimiterV1) Allow(now time.Time) (ok bool) {
	select {
	case <-limiter.exit:
		return false
	case limiter.ch <- struct{}{}: // always need to wait for a response which is not recommended
	}

	ok = limiter.allow(now)
	<-limiter.ch
	return
}

func (limiter *LimiterV1) AllowWithContext(ctx context.Context, now time.Time) (ok bool) {
	select {
	case <-limiter.exit:
		return false
	case <-ctx.Done(): // allow context like a timeout
		return false
	case limiter.ch <- struct{}{}:
	}

	ok = limiter.allow(now)
	<-limiter.ch
	return
}

func (limiter *LimiterV1) Last() time.Time {
	return limiter.vec[limiter.p]
}

func (limiter *LimiterV1) Stop() {
	close(limiter.exit)
}
