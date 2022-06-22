package rate_limit

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterV3_t1(t *testing.T) {
	limiter, _ := NewLimiterV3(time.Second, 3, false)

	task := func() {
		now := time.Now()
		fmt.Printf("~~~ %s get token: %t\n", now.Format(time.RFC3339), limiter.Allow(now))
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			go task()
		}
		time.Sleep(2 * time.Second)
		fmt.Println("")
	}
}

func TestLimiterV3_t2(t *testing.T) {
	limiter, _ := NewLimiterV3(time.Second, 10, false)

	task := func() {
		now := time.Now()
		fmt.Printf("~~~ %s get token: %t\n", now.Format(time.RFC3339), limiter.Allow(now))
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 20; j++ {
			go task()
		}
		time.Sleep(2 * time.Second)
		fmt.Println("")
	}
}

// go test  -run none  -bench ^BenchmarkLimiterV3_b1$ -count 5
// # 34 ns/op, sync.Mutex 88.72 ns/op,
func BenchmarkLimiterV3_b1(b *testing.B) {
	limiter, _ := NewLimiterV3(1*time.Second, 1, false)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}

// go test  -run none  -bench ^BenchmarkLimiterV3_b2$ -count 5
// # 63.45 ns/op
func BenchmarkLimiterV3_b2(b *testing.B) {
	limiter, _ := NewLimiterV3(time.Second, 1000, false)

	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(now)
		}
	})
}
