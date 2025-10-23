package testutil

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// NewTestRouter creates a Gin engine in test mode.
func NewTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// NewRecorder creates an HTTP recorder for testing.
func NewRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
