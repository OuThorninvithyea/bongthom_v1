package redis

import (

	// Commmnuity pacakges
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verify connection — crash if unavailable
	if err := Client.Ping(context.Background()).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	return Client
}
