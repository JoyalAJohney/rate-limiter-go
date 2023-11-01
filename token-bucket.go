package main

import (
	"time"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)


func TokenBucketRateLimit(ctx context.Context, redisClient *redis.Client, userId string,
	refillWindow int64, maximumRequests int64) bool {

		tokenBucketKey := "token_bucket:" + userId
		lastRefillTimeKey := "last_refill_time:" + userId

		lastRefillTimeInStr := redisClient.Get(ctx, lastRefillTimeKey).Val()
		lastRefillTime, _ := strconv.ParseInt(lastRefillTimeInStr, 10, 64)

		currentTime := time.Now().Unix()
		timeElapsed := currentTime - lastRefillTime
		
		// Calculate the number of tokens to be added to the bucket since the last refill
		if timeElapsed >= refillWindow {
			redisClient.Set(ctx, tokenBucketKey, strconv.FormatInt(maximumRequests, 10), 0)
			redisClient.Set(ctx, lastRefillTimeKey, strconv.FormatInt(currentTime, 10), 0)
		} else {
			token := redisClient.Get(ctx, tokenBucketKey).Val()
			tokenCount, _ := strconv.ParseInt(token, 10, 64)

			if tokenCount <= 0 {
				return false
			}
		}

		// Consume a token and proceed with request
		redisClient.Decr(ctx, tokenBucketKey)
		return true
}