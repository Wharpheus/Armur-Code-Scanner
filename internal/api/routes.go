package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/scan/repo", ScanHandler)
	r.POST("/api/v1/advanced-scan/repo", AdvancedScanResult)
	r.POST("/api/v1/scan/file", ScanFile)
	r.GET("/api/v1/status/:task_id", TaskStatus)
	r.GET("/api/v1/reports/owasp/:task_id", TaskOwasp)
	r.GET("/api/v1/reports/sans/:task_id", TaskSans)
}
