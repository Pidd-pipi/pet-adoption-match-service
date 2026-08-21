package middleware

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentRateLimitAndRequestCount(t *testing.T) {
	r := NewRateLimiter(100, time.Minute)
	const workers = 16
	const iters = 200
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			ip := fmt.Sprintf("10.0.0.%d", id%4)
			for j := 0; j < iters; j++ {
				r.Allow(ip)
				if j%20 == 0 {
					_ = r.Snapshot()
				}
				if j%50 == 0 {
					r.Reset(ip)
				}
				bumpRequestCount(ip)
			}
		}(w)
	}
	close(start)
	wg.Wait()
}
