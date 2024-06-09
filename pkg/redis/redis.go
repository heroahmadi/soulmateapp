package redis

import (
	"soulmateapp/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

func Set(key string, value string) error {
	rdb := config.GetRedisClient()
	ctx := config.GetContext()
	err := rdb.Set(ctx, key, value, 0).Err()

	return err
}

func Get(key string) (string, error) {
	rdb := config.GetRedisClient()
	ctx := config.GetContext()
	return rdb.Get(ctx, key).Result()
}

func SetHash(key string, field string, value string) error {
	rdb := config.GetRedisClient()
	ctx := config.GetContext()
	return rdb.HSet(ctx, key, field, value).Err()
}

func GetHash(key string, field string) (string, error) {
	rdb := config.GetRedisClient()
	ctx := config.GetContext()
	val, err := rdb.HGet(ctx, key, field).Result()
	if err != nil && err == redis.Nil {
		return "", nil
	}
	return val, err
}

func SetExpiryTime(key string, expiryTime time.Duration) error {
	rdb := config.GetRedisClient()
	ctx := config.GetContext()
	return rdb.Expire(ctx, key, expiryTime).Err()
}
