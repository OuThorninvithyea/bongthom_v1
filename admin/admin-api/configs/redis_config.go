package config

import (
	"os"

	"admin-api/pkg/utls"
)

type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisExpire   int
}

func InitRedis() *RedisConfig {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_db := utls.GetenvInt("REDIS_DB_NUMBER", 0)
	redis_expire := utls.GetenvInt("REDIS_EXPIRE", 60)

	return &RedisConfig{
		RedisHost:     redis_host,
		RedisPort:     redis_port,
		RedisPassword: redis_password,
		RedisDB:       redis_db,
		RedisExpire:   redis_expire,
	}
}
