package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"armur-codescanner/internal/redis"
)

var redisClient = redis.RedisClient()

func SaveTaskResult(taskID string, result map[string]interface{}) error {
	ctx := context.Background()

	resultData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, taskID, resultData, 24*time.Hour).Err()
}

func GetTaskResult(taskID string) (interface{}, error) {
	ctx := context.Background()

	resultData, err := redisClient.Get(ctx, taskID).Result()
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
