package main

import (
	"context"
	"time"
	"strconv"

	"github.com/go-redis/redis/v8"
)


func LeakyBucketThrottling(
	ctx context.Context, redisClient *redis.Client, 
	userId string, leakRate int64, maxBucketCapacity int64) bool {

	// Throttling so that the service handles only 5 requests per second
	bucketKey := "leaky_bucket"
	lastLeakTimeKey := "last_leak_time"

	currentTime := time.Now().Unix()
	lastLeakTimeInStr := redisClient.Get(ctx, lastLeakTimeKey).Val()
	lastLeakTime, _ := strconv.ParseInt(lastLeakTimeInStr, 10, 64)

	leakAmount := (currentTime - lastLeakTime) * leakRate

	// Perform the leak by trimming the queue
	redisClient.LTrim(ctx, bucketKey, leakAmount, -1)

	// update the last leak time
	redisClient.Set(ctx, lastLeakTimeKey, strconv.FormatInt(currentTime, 10), 0)

	// Check current size of bucket
	bucketLevel := redisClient.LLen(ctx, bucketKey).Val()

	if bucketLevel >= maxBucketCapacity {
		return false
	}

	// Add the current request to the bucket
	redisClient.RPush(ctx, bucketKey, strconv.FormatInt(currentTime, 10))
	return true
}