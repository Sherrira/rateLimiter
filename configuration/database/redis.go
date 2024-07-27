package database

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
}
