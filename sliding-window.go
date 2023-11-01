package main

import (
	"math/rand"
	"time"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)


func SlidingWindowRateLimit(ctx context.Context, redisClient *redis.Client, userId string,
	windowLengthInSeconds int64, maximumRequests int64) bool {

		currentTime := time.Now().Unix()
		windowStart := currentTime - windowLengthInSeconds

		sortedSetKey := "rate_limit:" + userId
		currentTimeInStr := strconv.FormatInt(currentTime, 10)
		windowStartInStr := strconv.FormatInt(windowStart, 10)
		
		// Count the number of requests in the current window
		requestCount := redisClient.ZCount(ctx, sortedSetKey, windowStartInStr, currentTimeInStr).Val()

		if requestCount >= maximumRequests {
			return false
		}

		/* 
			Add the current request to the sorted set, with the score as the current time
			Usually the member used can be a uniqueRequestId
		*/
		randomId := strconv.FormatInt(int64(rand.Intn(1000000)), 10)
		uniqueRequestId := currentTimeInStr + ":" + randomId
		redisClient.ZAdd(ctx, sortedSetKey, &redis.Z{
			Score: float64(currentTime),
			Member: uniqueRequestId,
		}).Result()
		
		
		/* 
			Optionally, remove expired requests from the sorted set (the ones that fall outside the window)
			ie: From -infinity to windowStart
		*/
		redisClient.ZRemRangeByScore(ctx, sortedSetKey, "-inf", windowStartInStr).Result()
		return true
}