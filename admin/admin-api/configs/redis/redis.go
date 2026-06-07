package redis

import (

	// Commmnuity pacakges
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	// Internal packages
	config "admin-api/configs"
)

var Client *redis.Client

func NewRedisClient() *redis.Client {
	cfg := config.InitRedis()

	host := cfg.RedisHost
	if host == "" {
		host = "localhost"
	}
	port := cfg.RedisPort
	if port == "" {
		port = "6379"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Verify connection — crash if unavailable
	if err := Client.Ping(context.Background()).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	return Client
}
