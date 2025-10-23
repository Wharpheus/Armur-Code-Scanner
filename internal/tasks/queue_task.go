package tasks

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func RedisClientOptions() asynq.RedisConnOpt {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := getEnvAsInt("REDIS_DB", 0)

	redisURL := "redis://"
	if password != "" {
		redisURL += ":" + password + "@"
	}
	redisURL += addr + "/" + strconv.Itoa(db)

	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		panic(err) // or handle accordingly, since this is config time
	}
	return opt
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

func EnqueueScanTask(scanType, repoURL, language string) (string, error) {
	taskID := uuid.New().String()

	taskPayload, err := json.Marshal(map[string]string{
		"repository_url": repoURL,
		"language":       language,
		"scan_type":      scanType,
		"task_id":        taskID,
	})
	if err != nil {
		return "", err
	}

	client := asynq.NewClient(RedisClientOptions())
	defer client.Close()

	task := asynq.NewTask("scan:repo", taskPayload)
	_, err = client.Enqueue(task, asynq.Queue("default"), asynq.MaxRetry(3), asynq.Timeout(30*time.Minute))
	if err != nil {
		return "", err
	}

	return taskID, nil
}

func EnqueueBatchScanTask(contractPaths []string, language, network string) (string, error) {
	taskID := uuid.New().String()

	taskPayload, err := json.Marshal(map[string]interface{}{
		"contract_paths": contractPaths,
		"language":       language,
		"network":        network,
		"scan_type":      "batch_scan",
		"task_id":        taskID,
	})
	if err != nil {
		return "", err
	}

	client := asynq.NewClient(RedisClientOptions())
	defer client.Close()

	task := asynq.NewTask("batch:scan", taskPayload)
	_, err = client.Enqueue(task, asynq.Queue("default"), asynq.MaxRetry(3), asynq.Timeout(60*time.Minute))
	if err != nil {
		return "", err
	}

	return taskID, nil
}
