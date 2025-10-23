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

	var taskData map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &taskData); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	repoURL, _ := taskData["repository_url"].(string)
	language, _ := taskData["language"].(string)
	scanType, _ := taskData["scan_type"].(string)
	taskID, _ := taskData["task_id"].(string)

	var result map[string]interface{}
	var err error
	switch scanType {
	case utils.SimpleScan:
		result = tasks.RunScanTask(repoURL, language)
	case utils.AdvancedScan:
		result = tasks.AdvancedScanRepositoryTask(repoURL, language)
	case utils.FileScan:
		result, _ = tasks.ScanFileTask(repoURL)
	case utils.LocalScan:
		result = tasks.RunScanTaskLocal(repoURL, language)
	case "batch_scan":
		result = tasks.RunBatchScanTask(taskData)
	default:
		return fmt.Errorf("unknown scan type: %s", scanType)
	}

	err = tasks.SaveTaskResult(taskID, result)
	if err != nil {
		return fmt.Errorf("failed to store scan result: %w", err)
	}

	return nil
}
