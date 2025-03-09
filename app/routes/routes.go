package routes

import (
	"net/http"

	handleFull "ratelimit/app/handlers"
	"ratelimit/app/middleware"
)

func SetupRoutes() {
	http.Handle("/full", middleware.RateLimit(http.HandlerFunc(handleFull.HandleFull)))
	http.HandleFunc("/generate-token", handleFull.GenerateToken) // Nova rota para gerar o token
}
