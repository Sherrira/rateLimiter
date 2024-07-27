package middleware_redis

import (
	"context"
	"log"
	"time"

	"github.com/Sherrira/rateLimiter/configuration/database"
	"github.com/Sherrira/rateLimiter/internal/infra/api/middleware"
	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	client *redis.Client
	config middleware.RateLimiterConfiguration
}

func NewRateLimiter(config middleware.RateLimiterConfiguration) middleware.RateLimiter {
	client := database.NewRedisClient()
	return &RedisRateLimiter{
		client: client,
		config: config,
	}
}

func (rl *RedisRateLimiter) IsRateLimited(ctx context.Context, key string, limit int) bool {
	blocked, err := rl.client.Get(ctx, "block:"+key).Result()
	if err == nil && blocked == "1" {
		log.Printf("Key %s is blocked", key)
		return true
	}

	requests, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Error incrementing key %s: %v", key, err)
		return true
	}

	if requests == 1 {
		rl.client.Expire(ctx, key, time.Second)
	}

	if requests > int64(limit) {
		log.Printf("Rate limit exceeded for key %s", key)
		rl.client.SetEX(ctx, "block:"+key, 1, rl.config.BlockDuration)
		return true
	}

	return false
}
