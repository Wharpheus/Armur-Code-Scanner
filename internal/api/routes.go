package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Scan routes
		api.POST("/scan/repo", ScanHandler)
		api.POST("/advanced-scan/repo", AdvancedScanResult)
		api.POST("/scan/file", ScanFile)
		api.POST("/scan/local", ScanLocalHandler)

		// status
		api.GET("/status/:task_id", TaskStatus)

		// reports
		api.GET("/reports/owasp/:task_id", TaskOwasp)
		api.GET("/reports/sans/:task_id", TaskSans)
	}

}
