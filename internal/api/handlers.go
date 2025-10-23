package api

import (
	"armur-codescanner/internal/tasks"
	utils "armur-codescanner/pkg"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ScanRequest represents a scan request for a github repository with a specified language
type ScanRequest struct {
	RepositoryURL string `json:"repository_url" example:"https://github.com/Armur-Ai/Armur-Code-Scanner"`
	Language      string `json:"language" example:"go"`
}

// LocalScanRequest represents a scan request for a local repository with a specified language
type LocalScanRequest struct {
	LocalPath string `json:"local_path" binding:"required" example:"/armur/repo"`
	Language  string `json:"language" example:"go"`
}

// BatchScanRequest represents a batch scan request for multiple contracts
type BatchScanRequest struct {
	ContractPaths []string `json:"contract_paths" binding:"required" example:"[\"/armur/contracts/Token.sol\", \"/armur/contracts/LP.sol\"]"`
	Language      string   `json:"language" example:"solidity"`
	Network       string   `json:"network" example:"polygon"`
}

// ScanHandler godoc
// @Summary Trigger a code scan on a repository.
// @Description Enqueues a scan task for a given repository URL and language.
// @Tags scan
// @Accept json
// @Produce json
// @Param request body ScanRequest true "Request body containing repository URL and language"
// @Success 200 {object} map[string]string  "Successfully enqueued task"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Failed to enqueue scan task"
// @Router /api/v1/scan/repo [post]
func ScanHandler(c *gin.Context) {
	var request ScanRequest

	// Bind JSON to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate repository URL and language
	if request.RepositoryURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository URL is required"})
		return
	}
	if request.Language != "" && request.Language != "go" && request.Language != "py" && request.Language != "js" && request.Language != "solidity" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language"})
		return
	}

	// Enqueue the scan task
	taskID, err := tasks.EnqueueScanTask(utils.SimpleScan, request.RepositoryURL, request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue scan task", "details": err.Error()})
		return
	}

	// Respond with the Task ID
	c.JSON(http.StatusOK, gin.H{"task_id": taskID})
}

// AdvancedScanResult godoc
// @Summary Trigger a code scan on a repository with advanced scans.
// @Description Enqueues an advanced scan task for a given repository URL and language.
// @Tags scan
// @Accept json
// @Produce json
// @Param request body ScanRequest true "Request body containing repository URL and language"
// @Success 200 {object} map[string]string "Successfully enqueued task"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Failed to enqueue scan task"
// @Router /api/v1/advanced-scan/repo [post]
func AdvancedScanResult(c *gin.Context) {
	var request ScanRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.RepositoryURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository URL is required"})
		return
	}
	if request.Language != "" && request.Language != "go" && request.Language != "py" && request.Language != "js" && request.Language != "solidity" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language"})
		return
	}

	taskID, err := tasks.EnqueueScanTask(utils.AdvancedScan, request.RepositoryURL, request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue scan task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_id": taskID})
}

// ScanFile godoc
// @Summary Trigger a code scan on file.
// @Description Enqueues a scan task for a given file.
// @Tags scan
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to be scanned"
// @Success 202 {object} map[string]string "Successfully enqueued task"
// @Failure 400 {object} map[string]string "No file part or no selected file"
// @Failure 500 {object} map[string]string "Failed to create temp directory or Failed to save file"
// @Router /api/v1/scan/file [post]
func ScanFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file part or no selected file"})
		return
	}

	baseDir := "/tmp/armur/repos"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base directory", "details": err.Error()})
			return
		}
	}

	tempDir, err := os.MkdirTemp(baseDir, "scan")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory", "details": err.Error()})
		return
	}

	filePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}

	taskID, err := tasks.EnqueueScanTask(utils.FileScan, filePath, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue scan task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"task_id": taskID,
	})
}

// TaskStatus godoc
// @Summary Get task status and results.
// @Description Get the status and results of a scan task by its ID.
// @Tags scan
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]interface{} "Successfully retrieved task result"
// @Router /api/v1/status/{task_id} [get]
func TaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")

	result, err := tasks.GetTaskResult(taskID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "pending",
			"task_id": taskID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"task_id": taskID,
	})
}

// TaskOwasp godoc
// @Summary Get OWASP report for a task result.
// @Description Generates OWASP report from a specific task result using task ID.
// @Tags report
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} []utils.ReportItem "Successfully generated OWASP report"
// @Failure 404 {object} map[string]string "Task result not found"
// @Failure 500 {object} map[string]string "Failed to fetch task result"
// @Router /api/v1/reports/owasp/{task_id} [get]
func TaskOwasp(c *gin.Context) {
	taskID := c.Param("task_id")

	// Fetch task result from Redis
	taskResult, err := tasks.GetTaskResult(taskID)
	fmt.Println("taskresult ", taskResult)
	if err != nil {
		if err.Error() == "task result not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task result not found. Pls wait"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task result", "details": err.Error()})
		}
		return
	}

	report, err := utils.GenerateOwaspReport(taskResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// TaskSans godoc
// @Summary Get SANS report for a task result.
// @Description Generates SANS report from a specific task result using task ID.
// @Tags report
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} []utils.SANSReportItem "Successfully generated SANS report"
// @Failure 404 {object} map[string]string "Task result not found"
// @Failure 500 {object} map[string]string "Failed to fetch task result"
// @Router /api/v1/reports/sans/{task_id} [get]
func TaskSans(c *gin.Context) {
	taskID := c.Param("task_id")

	// Fetch task result from Redis
	taskResult, err := tasks.GetTaskResult(taskID)
	fmt.Println("taskresult ", taskResult)
	if err != nil {
		if err.Error() == "task result not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task result not found. Pls wait"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task result", "details": err.Error()})
		}
		return
	}

	report, err := utils.GenerateSANSReports(taskResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// ScanLocalHandler godoc
// @Summary Trigger a code scan on a local repository.
// @Description Enqueues a scan task for a given local repository path and language.
// @Tags scan
// @Accept json
// @Produce json
// @Param request body LocalScanRequest true "Request body containing local path and language"
// @Success 200 {object} map[string]string "Successfully enqueued task"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Failed to enqueue scan task"
// @Router /api/v1/scan/local [post]
func ScanLocalHandler(c *gin.Context) {
	var request LocalScanRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.LocalPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Local path is required"})
		return
	}

	if request.Language != "" && request.Language != "go" && request.Language != "py" && request.Language != "js" && request.Language != "solidity" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language"})
		return
	}

	// Verify the path exists and is a directory
	if info, err := os.Stat(request.LocalPath); os.IsNotExist(err) || !info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-existent directory"})
		return
	}

	// Enqueue the scan task with "local" scan type
	taskID, err := tasks.EnqueueScanTask(utils.LocalScan, request.LocalPath, request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue scan task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_id": taskID})
}

// BatchScanHandler godoc
// @Summary Trigger batch gas optimization and LP pairing scans for multiple contracts.
// @Description Enqueues batch scan tasks for multiple Solidity contracts with gas optimization and DeFi analysis.
// @Tags scan
// @Accept json
// @Produce json
// @Param request body BatchScanRequest true "Request body containing contract paths and network"
// @Success 200 {object} map[string]string "Successfully enqueued batch scan task"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Failed to enqueue batch scan task"
// @Router /api/v1/batch-scan/contracts [post]
func BatchScanHandler(c *gin.Context) {
	var request BatchScanRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.ContractPaths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract paths are required"})
		return
	}

	if request.Language != "solidity" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch scanning currently only supports Solidity"})
		return
	}

	if request.Network != "polygon" && request.Network != "amoy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Network must be 'polygon' or 'amoy'"})
		return
	}

	// Validate all contract paths exist
	for _, path := range request.ContractPaths {
		if info, err := os.Stat(path); os.IsNotExist(err) || info.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid contract path: %s", path)})
			return
		}
	}

	// Enqueue batch scan task
	taskID, err := tasks.EnqueueBatchScanTask(request.ContractPaths, request.Language, request.Network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue batch scan task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": fmt.Sprintf("Batch scanning %d contracts for %s network", len(request.ContractPaths), request.Network),
	})
}
