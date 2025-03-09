package limiter

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter struct {
	storage       Storage
	requests      int
	timeWindow    time.Duration
	blockDuration time.Duration
}

func NewRateLimiter(storage Storage, requests int, timeWindow, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		storage:       storage,
		requests:      requests,
		timeWindow:    timeWindow,
		blockDuration: blockDuration,
	}
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	blocked, err := rl.storage.Get(ctx, blockKey)
	if err == nil && blocked == "1" {
		return false, nil
	}

	requestKey := fmt.Sprintf("req:%s", key)
	requestCount, err := rl.storage.Incr(ctx, requestKey)
	if err != nil {
		return false, err
	}

	if err := rl.storage.Expire(ctx, requestKey, rl.timeWindow); err != nil {
		return false, err
	}

	if requestCount > int64(rl.requests) {
		if err := rl.storage.Set(ctx, blockKey, "1", rl.blockDuration); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
