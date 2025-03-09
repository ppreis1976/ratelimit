package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() *miniredis.Miniredis {
	os.Setenv("REQUESTS_PER_SECOND", "5")
	os.Setenv("TIME_WINDOW", "1s")
	os.Setenv("BLOCK_DURATION", "10s")

	// Start a miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	InitRateLimiter(10, time.Second, 10*time.Second, client)

	return mr
}

func TestRateLimitWithoutToken(t *testing.T) {
	mr := setup()
	defer mr.Close()

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimitWithToken(t *testing.T) {
	mr := setup()
	defer mr.Close()

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	mr.Set("test_token", "10")

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	req.Header.Set("API_KEY", "Bearer test_token")

	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimitInvalidToken(t *testing.T) {
	mr := setup()
	defer mr.Close()

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	req.Header.Set("API_KEY", "Bearer invalid_token")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
