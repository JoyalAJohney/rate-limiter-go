package main

import (
	"fmt"
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
}

func main() {

	// Configurations
	userId := "user123"
	// leakRate := int64(5) // 5 request per second
	refillWindow := int64(60) // 60 seconds
	// windowLengthInSeconds := int64(60) // window size is 60 seconds
	maximumTokens := int64(10) // maximum 10 requests per window

	for i := 1; i <= 25; i++ {
		isAllowed := TokenBucketRateLimit(ctx, redisClient, userId, refillWindow, maximumTokens)
		fmt.Printf("Request %d status - allowed: %t \n", i, isAllowed)
	}

	http.ListenAndServe(":8080", nil)
}