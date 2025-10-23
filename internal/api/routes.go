package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// securityHeadersMiddleware adds standard security headers to every response.
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		// Set HSTS only when behind HTTPS (respect X-Forwarded-Proto)
		if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto == "https" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		c.Next()

		// Ensure content type is set for JSON by default when status has body
		if c.Writer.Header().Get("Content-Type") == "" && c.Writer.Status() >= 200 && c.Writer.Status() != http.StatusNoContent {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
	}
}

func RegisterRoutes(r *gin.Engine) {
	// apply security headers globally
	r.Use(securityHeadersMiddleware())

	api := r.Group("/api/v1")
	{
		// Scan routes
		api.POST("/scan/repo", ScanHandler)
		api.POST("/advanced-scan/repo", AdvancedScanResult)
		api.POST("/scan/file", ScanFile)
		api.POST("/scan/local", ScanLocalHandler)

		// Batch scan route for multiple contracts
		api.POST("/batch-scan/contracts", BatchScanHandler)

		// status
		api.GET("/status/:task_id", TaskStatus)

		// reports
		api.GET("/reports/owasp/:task_id", TaskOwasp)
		api.GET("/reports/sans/:task_id", TaskSans)
	}

}
