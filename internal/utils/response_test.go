package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := gin.H{"test": "data"}
	Success(c, testData)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify all required fields are present
	assert.Equal(t, CodeSuccess, response.Code)
	assert.Equal(t, "success", response.Message)
	assert.NotNil(t, response.Data)
	assert.NotZero(t, response.Timestamp)
	assert.Empty(t, response.Error)
}

func TestSuccessWithMessageResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	customMessage := "operation completed"
	testData := gin.H{"result": "ok"}
	SuccessWithMessage(c, customMessage, testData)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeSuccess, response.Code)
	assert.Equal(t, customMessage, response.Message)
	assert.NotNil(t, response.Data)
	assert.NotZero(t, response.Timestamp)
}

func TestSuccessPaginatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := []gin.H{{"id": 1}, {"id": 2}}
	pagination := &Pagination{
		Page:       1,
		PageSize:   20,
		Total:      100,
		TotalPages: 5,
	}

	SuccessPaginated(c, testData, pagination)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify all required fields are present
	assert.Equal(t, CodeSuccess, response.Code)
	assert.Equal(t, "success", response.Message)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Pagination)
	assert.NotZero(t, response.Timestamp)

	// Verify pagination fields
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 20, response.Pagination.PageSize)
	assert.Equal(t, 100, response.Pagination.Total)
	assert.Equal(t, 5, response.Pagination.TotalPages)
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	appErr := NewAppError(CodeInvalidParams, "invalid parameters", nil)
	Error(c, appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify all required fields are present
	assert.Equal(t, CodeInvalidParams, response.Code)
	assert.Equal(t, "invalid parameters", response.Message)
	assert.NotZero(t, response.Timestamp)
	assert.Nil(t, response.Data)
}

func TestErrorResponseWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	underlyingErr := assert.AnError
	appErr := NewAppError(CodeDatabaseError, "database error", underlyingErr)
	Error(c, appErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify all required fields are present
	assert.Equal(t, CodeDatabaseError, response.Code)
	assert.Equal(t, "database error", response.Message)
	assert.NotZero(t, response.Timestamp)
	assert.NotEmpty(t, response.Error)
}

func TestErrorWithMessageResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorWithMessage(c, CodeNotFound, "resource not found")

	// CodeNotFound (40401) maps to 404 Not Found
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeNotFound, response.Code)
	assert.Equal(t, "resource not found", response.Message)
	assert.NotZero(t, response.Timestamp)
}

func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		pageSize  int
		total     int
		wantPage  int
		wantSize  int
		wantTotal int
		wantPages int
	}{
		{
			name:      "normal pagination",
			page:      2,
			pageSize:  20,
			total:     100,
			wantPage:  2,
			wantSize:  20,
			wantTotal: 100,
			wantPages: 5,
		},
		{
			name:      "page less than 1",
			page:      0,
			pageSize:  20,
			total:     100,
			wantPage:  1,
			wantSize:  20,
			wantTotal: 100,
			wantPages: 5,
		},
		{
			name:      "page size less than 1",
			page:      1,
			pageSize:  0,
			total:     100,
			wantPage:  1,
			wantSize:  10,
			wantTotal: 100,
			wantPages: 10,
		},
		{
			name:      "partial last page",
			page:      1,
			pageSize:  20,
			total:     95,
			wantPage:  1,
			wantSize:  20,
			wantTotal: 95,
			wantPages: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := CalculatePagination(tt.page, tt.pageSize, tt.total)
			assert.Equal(t, tt.wantPage, pagination.Page)
			assert.Equal(t, tt.wantSize, pagination.PageSize)
			assert.Equal(t, tt.wantTotal, pagination.Total)
			assert.Equal(t, tt.wantPages, pagination.TotalPages)
		})
	}
}

func TestHTTPStatusCodeMapping(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		wantStatus int
	}{
		{"success", CodeSuccess, http.StatusOK},
		{"invalid params", CodeInvalidParams, http.StatusBadRequest},
		{"unauthorized", CodeUnauthorized, http.StatusUnauthorized},
		{"not found", CodeNotFound, http.StatusNotFound},
		{"too many requests", CodeTooManyRequests, http.StatusTooManyRequests},
		{"internal error", CodeInternalError, http.StatusInternalServerError},
		{"database error", CodeDatabaseError, http.StatusInternalServerError},
		{"AI service error", CodeAIServiceError, http.StatusInternalServerError},
		{"unknown 4xxxx code", 40999, http.StatusBadRequest},
		{"unknown 5xxxx code", 59999, http.StatusInternalServerError},
		{"completely unknown code", 99999, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := getHTTPStatusCode(tt.code)
			assert.Equal(t, tt.wantStatus, status)
		})
	}
}

func TestResponseTimestamp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	beforeTime := time.Now().Unix()
	Success(c, gin.H{"test": "data"})
	afterTime := time.Now().Unix()

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify timestamp is within reasonable range
	assert.GreaterOrEqual(t, response.Timestamp, beforeTime)
	assert.LessOrEqual(t, response.Timestamp, afterTime)
}

func TestSanitizeErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sanitize DSN password",
			input:    "user:password123@tcp(localhost:3306)/dbname?charset=utf8mb4",
			expected: "user:***@tcp(localhost:3306)/dbname?charset=utf8mb4",
		},
		{
			name:     "no sensitive data",
			input:    "simple error message",
			expected: "simple error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeErrorMessage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
