package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedisClient(client *redis.Client) {
	redisClient = client
}

type TokenRequest struct {
	RequestsPerSecond int `json:"requests_per_second"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token := generateRandomToken()
	err := redisClient.Set(context.Background(), token, req.RequestsPerSecond, 0).Err()
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateRandomToken() string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
