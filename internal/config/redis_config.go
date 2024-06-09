package config

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	clientInstance *redis.Client
	clientOnce     sync.Once
)

func initRedisClient() {
	clientInstance = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	clientOnce.Do(initRedisClient)
	return clientInstance
}

func GetContext() context.Context {
	return context.Background()
}
