package redis

import (
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

var (
	client     *redis.Client
	clientOnce sync.Once
)

func RedisClientOptions() asynq.RedisConnOpt {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := getEnvAsInt("REDIS_DB", 0)

	return asynq.RedisClientOpt{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}

func RedisClient() *redis.Client {
	clientOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		password := os.Getenv("REDIS_PASSWORD")
		db := getEnvAsInt("REDIS_DB", 0)

		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	})
	return client
}

func getEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
