package utils

import (
	"sync"
	"time"
)

type RateLimiter struct {
	requests    []time.Time
	maxRequests int
	window      time.Duration
	mutex       sync.Mutex
}

var secRateLimiter = &RateLimiter{
	requests:    make([]time.Time, 0),
	maxRequests: 4,
	window:      time.Second,
}

func (rl *RateLimiter) Wait() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	cutoff := now.Add(-rl.window)
	validRequests := make([]time.Time, 0)
	for _, req := range rl.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	rl.requests = validRequests

	if len(rl.requests) >= rl.maxRequests {
		oldestRequest := rl.requests[0]
		waitTime := oldestRequest.Add(rl.window).Sub(now)
		if waitTime > 0 {
			time.Sleep(waitTime)
		}
	}

	rl.requests = append(rl.requests, time.Now())
}
