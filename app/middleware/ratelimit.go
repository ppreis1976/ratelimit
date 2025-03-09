package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient          *redis.Client
	requestsPerSecond    int
	timeWindow           time.Duration
	blockDuration        time.Duration
	defaultRPS           int
	defaultTimeWindow    time.Duration
	defaultBlockDuration time.Duration
)

func InitRateLimiter(rps int, tw, bd time.Duration, client *redis.Client) {
	requestsPerSecond = rps
	timeWindow = tw
	blockDuration = bd
	redisClient = client

	defaultRPS, _ = strconv.Atoi(os.Getenv("REQUESTS_PER_SECOND"))
	defaultTimeWindow, _ = time.ParseDuration(os.Getenv("TIME_WINDOW"))
	defaultBlockDuration, _ = time.ParseDuration(os.Getenv("BLOCK_DURATION"))
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		var rps int
		var tw, bd time.Duration

		if token == "" {
			// Use default rate limit for requests without a token
			rps = defaultRPS
			tw = defaultTimeWindow
			bd = defaultBlockDuration
		} else {
			token = token[len("Bearer "):]
			requestsPerSecondStr, err := redisClient.Get(context.Background(), token).Result()
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			rps, err = strconv.Atoi(requestsPerSecondStr)
			if err != nil {
				http.Error(w, "Invalid rate limit value", http.StatusInternalServerError)
				return
			}
			tw = timeWindow
			bd = blockDuration
		}

		// Implement rate limiting logic here using rps, tw, and bd
		key := "rate_limit:" + token
		currentTime := time.Now().Unix()
		allowedRequests := int64(rps)
		bucketSize := allowedRequests

		// Get the current state of the bucket
		bucket, err := redisClient.HGetAll(context.Background(), key).Result()
		if err != nil && err != redis.Nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		lastTime, _ := strconv.ParseInt(bucket["last_time"], 10, 64)
		tokens, _ := strconv.ParseInt(bucket["tokens"], 10, 64)

		// Calculate the number of tokens to add based on the elapsed time
		elapsedTime := currentTime - lastTime
		newTokens := elapsedTime * allowedRequests / int64(tw.Seconds())
		tokens = min(bucketSize, tokens+newTokens)

		if tokens > 0 {
			tokens--
			redisClient.HSet(context.Background(), key, map[string]interface{}{
				"last_time": currentTime,
				"tokens":    tokens,
			})
			next.ServeHTTP(w, r)
		} else {
			redisClient.Set(context.Background(), "block:"+token, "1", bd)
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		}
	})
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
