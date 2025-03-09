package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ratelimit/app/handlers"
	"ratelimit/app/middleware"
	"ratelimit/app/routes"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initRateLimiter() {
	requestsPerSecond, err := strconv.Atoi(os.Getenv("REQUESTS_PER_SECOND"))
	if err != nil {
		log.Fatalf("Invalid REQUESTS_PER_SECOND value")
	}

	timeWindowStr := os.Getenv("TIME_WINDOW")
	var timeWindow time.Duration
	if strings.HasSuffix(timeWindowStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(timeWindowStr, "m"))
		if err != nil {
			log.Fatalf("Invalid TIME_WINDOW value")
		}
		timeWindow = time.Duration(minutes) * time.Minute
	} else if strings.HasSuffix(timeWindowStr, "s") {
		seconds, err := strconv.Atoi(strings.TrimSuffix(timeWindowStr, "s"))
		if err != nil {
			log.Fatalf("Invalid TIME_WINDOW value")
		}
		timeWindow = time.Duration(seconds) * time.Second
	} else {
		log.Fatalf("Invalid TIME_WINDOW format")
	}

	blockDurationStr := os.Getenv("BLOCK_DURATION")
	var blockDuration time.Duration
	if strings.HasSuffix(blockDurationStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(blockDurationStr, "m"))
		if err != nil {
			log.Fatalf("Invalid BLOCK_DURATION value")
		}
		blockDuration = time.Duration(minutes) * time.Minute
	} else if strings.HasSuffix(blockDurationStr, "s") {
		seconds, err := strconv.Atoi(strings.TrimSuffix(blockDurationStr, "s"))
		if err != nil {
			log.Fatalf("Invalid BLOCK_DURATION value")
		}
		blockDuration = time.Duration(seconds) * time.Second
	} else {
		log.Fatalf("Invalid BLOCK_DURATION format")
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	middleware.InitRateLimiter(requestsPerSecond, timeWindow, blockDuration, client)
	handlers.InitRedisClient(client)
}

func startServer() {
	log.Println("Starting server on :8080")
	routes.SetupRoutes()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	loadEnv()
	initRateLimiter()
	startServer()
}
