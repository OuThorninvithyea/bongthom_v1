package redis_util

import (

	// Community packages
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisUtil(client *redis.Client, ttl time.Duration) *RedisUtil {
	return &RedisUtil{client: client, ttl: ttl}
}

func (r *RedisUtil) SetCacheKey(key string, value interface{}, ctx context.Context) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *RedisUtil) GetCacheKey(key string, dest interface{}, ctx context.Context) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}
