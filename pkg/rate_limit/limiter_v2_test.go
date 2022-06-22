package rate_limit

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterV2_t1(t *testing.T) {
	limiter, _ := NewLimiterV2(10*time.Second, 3)

	for i := 0; i < 20; i++ {
		now := time.Now()
		fmt.Printf("~~~ %s, %t\n", now.Format(time.RFC3339), limiter.Allow(now))
		time.Sleep(time.Second)
	}
}

// go test  -run none  -bench ^BenchmarkLimiterV2_b1$ -count 5
// # 312 ns/op
func BenchmarkLimiterV2_b1(b *testing.B) {
	limiter, _ := NewLimiterV2(1*time.Second, 1)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}

// go test  -run none  -bench ^BenchmarkLimiterV2_b2$ -count 5
// # 313 ns/op
func BenchmarkLimiterV2_b2(b *testing.B) {
	limiter, _ := NewLimiterV2(time.Second, 1000)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}
