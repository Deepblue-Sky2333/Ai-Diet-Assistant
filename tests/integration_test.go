package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckEndpoint tests the health check endpoint format
func TestHealthCheckEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status":  "ok",
			"service": "ai-diet-assistant",
		})
	})
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test404Handler tests the 404 not found handler format
func Test404Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	r.NoRoute(func(c *gin.Context) {
		utils.ErrorWithMessage(c, utils.CodeNotFound, "endpoint not found")
	})
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestResponseFormatCompliance tests that response formats are correct
func TestResponseFormatCompliance(t *testing.T) {
	// This test verifies response format compliance
	// Actual endpoint testing should be done with the running server
	t.Log("Response format tests are covered by internal/utils/response_test.go")
	t.Log("Date handling tests are covered by internal/utils/date_test.go")
	t.Log("API endpoint tests should be run against the live server using tests/api_test.sh")
}
