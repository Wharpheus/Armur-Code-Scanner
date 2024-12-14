package worker

import (
	"armur-codescanner/internal/tasks"
	utils "armur-codescanner/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

type ScanTaskHandler struct{}

func (h *ScanTaskHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {

	var taskData map[string]string
	if err := json.Unmarshal(task.Payload(), &taskData); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	repoURL := taskData["repository_url"]
	language := taskData["language"]
	scanType := taskData["scan_type"]
	taskID := taskData["task_id"]

	var result map[string]interface{}
	var err error
	switch scanType {
	case utils.SimpleScan:
		result = tasks.RunScanTask(repoURL, language)
	case utils.AdvancedScan:
		result = tasks.AdvancedScanRepositoryTask(repoURL, language)
	case utils.FileScan:
		result, _ = tasks.ScanFileTask(repoURL)
	default:
		return fmt.Errorf("unknown scan type: %s", scanType)
	}

	err = tasks.SaveTaskResult(taskID, result)
	if err != nil {
		return fmt.Errorf("failed to store scan result: %w", err)
	}

	return nil
}
