package rate_limit

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterV1_t1(t *testing.T) {
	limiter, _ := NewLimiterV1(10*time.Second, 3, false)

	for i := 0; i < 20; i++ {
		now := time.Now()
		fmt.Printf("~~~ %s, %t\n", rfc3339now(), limiter.Allow(now))
		time.Sleep(time.Second)
	}
}

func TestLimiterV1_t2(t *testing.T) {
	limiter, _ := NewLimiterV1(10*time.Second, 3, true)

	for i := 0; i < 20; i++ {
		now := time.Now()
		fmt.Printf("~~~ %s, %t\n", rfc3339now(), limiter.Allow(now))
		time.Sleep(time.Second)
	}
}

// go test  -run none  -bench ^BenchmarkLimiterV1_b1$ -count 5
// # 1121 ns/op
func BenchmarkLimiterV1_b1(b *testing.B) {
	limiter, _ := NewLimiterV1(1*time.Second, 1, false)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}

// go test  -run none  -bench ^BenchmarkLimiterV1_b2$ -count 5
// # 1111 ns/op
func BenchmarkLimiterV1_b2(b *testing.B) {
	limiter, _ := NewLimiterV1(time.Second, 1000, false)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}
