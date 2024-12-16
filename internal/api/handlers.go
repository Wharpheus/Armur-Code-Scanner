package api

import (
	"armur-codescanner/internal/tasks"
	utils "armur-codescanner/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

type ScanRequest struct {
	RepositoryURL string `json:"repository_url"`
	Language      string `json:"language"`
}

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
	if request.Language != "" && request.Language != "go" && request.Language != "py" && request.Language != "js" {
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
	if request.Language != "" && request.Language != "go" && request.Language != "py" && request.Language != "js" {
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

func ScanFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file part or no selected file"})
		return
	}

	baseDir := "/armur/repos"
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

	//log.Printf("File path: %s", filePath)

	taskID, err := tasks.EnqueueScanTask(utils.FileScan, filePath, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue scan task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"task_id":    taskID,
	})
}

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