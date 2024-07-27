package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sherrira/rateLimiter/internal/infra/api/auth"
	"github.com/Sherrira/rateLimiter/internal/infra/api/handler"
	"github.com/Sherrira/rateLimiter/internal/infra/api/middleware"
	"github.com/Sherrira/rateLimiter/internal/infra/api/middleware/middleware_redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ipRateLimit, _ := strconv.Atoi(os.Getenv("IP_RATE_LIMIT"))
	blockDuration, _ := time.ParseDuration(os.Getenv("BLOCK_DURATION"))

	tokenLimitsEnv := os.Getenv("TOKEN_LIMITS")
	tokenLimits, err := parseTokenLimits(tokenLimitsEnv)
	if err != nil {
		log.Fatalf("Error parsing token limits: %v", err)
	}

	r := mux.NewRouter()

	rlConfig := middleware.RateLimiterConfiguration{
		IpRateLimit:   ipRateLimit,
		TokenLimits:   tokenLimits,
		BlockDuration: blockDuration,
	}
	limiter := middleware_redis.NewRateLimiter(rlConfig)
	authorizer := auth.NewAuthorizer()
	r.Use(middleware.RateLimitMiddleware(limiter, ipRateLimit, authorizer))

	r.HandleFunc("/", handler.Hello).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func parseTokenLimits(envVar string) (map[string]int, error) {
	tokenLimits := make(map[string]int)
	pairs := strings.Split(envVar, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid token limit pair: %s", pair)
		}
		key := kv[0]
		value, err := strconv.Atoi(kv[1])
		if err != nil {
			return nil, fmt.Errorf("invalid token limit value for key %s: %s", key, kv[1])
		}
		tokenLimits[key] = value
	}
	return tokenLimits, nil
}
