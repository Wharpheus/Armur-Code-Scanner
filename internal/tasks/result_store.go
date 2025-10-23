package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	resultStoreRedisClient     *redis.Client
	resultStoreRedisClientOnce sync.Once
)

func initResultStoreRedisClient() *redis.Client {
	resultStoreRedisClientOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		password := os.Getenv("REDIS_PASSWORD")
		db := getEnvAsIntResultStore("REDIS_DB", 0)

		resultStoreRedisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	})
	return resultStoreRedisClient
}

func getEnvAsIntResultStore(key string, defaultValue int) int {
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

func SaveTaskResult(taskID string, result map[string]any) error {
	ctx := context.Background()

	resultData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	client := initResultStoreRedisClient()
	return client.Set(ctx, taskID, resultData, 24*time.Hour).Err()
}

func GetTaskResult(taskID string) (any, error) {
	ctx := context.Background()

	client := initResultStoreRedisClient()
	resultData, err := client.Get(ctx, taskID).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, errors.New("task result not found")
		}
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(resultData), &result); err != nil {
		return nil, err
	}
	return result, nil
}
