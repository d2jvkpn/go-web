package rate_limit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext_t1(t *testing.T) {
	ctx := context.TODO()
	fmt.Println(ctx.Done() == nil) // true

	<-ctx.Done() // aways block
	fmt.Println("done!")
}

func TestContext_t2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	fmt.Println("context.Done() is nil:", ctx.Done() == nil) // false
	defer cancel()

	for i := 0; i < 8; i++ {
		select {
		case <-ctx.Done():
			fmt.Println(rfc3339now(), "time is out!")
			return
		default:
			fmt.Println(rfc3339now(), "continue your work")
		}
		time.Sleep(time.Second)
	}
}

func TestChannel(t *testing.T) {
	c1 := make(chan struct{})
	fmt.Println(c1 == nil)
	close(c1)
	fmt.Println(c1 == nil)

	c2 := make(chan bool, 1)
	fmt.Println(c2 == nil)
	close(c2)
	fmt.Println(c2 == nil)
}
