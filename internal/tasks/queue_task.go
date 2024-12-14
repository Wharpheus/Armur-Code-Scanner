package tasks

import (
	"armur-codescanner/internal/redis"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"time"
)

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

	client := asynq.NewClient(redis.RedisClientOptions())
	defer client.Close()

	task := asynq.NewTask("scan:repo", taskPayload)
	_, err = client.Enqueue(task, asynq.Queue("default"), asynq.MaxRetry(3), asynq.Timeout(30*time.Minute))
	if err != nil {
		return "", err
	}

	return taskID, nil
}
