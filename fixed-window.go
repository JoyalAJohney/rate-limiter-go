package main

import (
	"context"
	"time"
	"strconv"

	"github.com/go-redis/redis/v8"
)


func FixedWindowRateLimit(ctx context.Context, redisClient *redis.Client, userId string,
	windowLengthInSeconds int64, maximumRequests int64) bool {
	/* 
		Gives me number of windows till now,
		Ex: 2343 windows till now means, the current window is 2343
	*/
	currentWindow := strconv.FormatInt(time.Now().Unix() / windowLengthInSeconds, 10)
	
	// Performing rate limiting per user
	key := userId + ":" + currentWindow

	val := redisClient.Incr(ctx, key).Val()

	// If the value is 1, then it is the first request in the current window, so set the TTL
	if val == 1 {
		ttl := time.Duration(windowLengthInSeconds) * time.Second
		redisClient.Expire(ctx, key, ttl)
	}

	if val > maximumRequests {
		return false
	}
	return true
}