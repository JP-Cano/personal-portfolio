package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestValidateRequest_ValidExperience(t *testing.T) {
	tests := []struct {
		name        string
		requestBody dto.ExperienceRequest
		wantErr     bool
	}{
		{
			name: "Valid experience with all fields",
			requestBody: dto.ExperienceRequest{
				Title:       "Senior Software Engineer",
				Company:     "Tech Corp",
				URL:         stringPtr("https://techcorp.com"),
				Location:    "San Francisco, CA",
				Type:        "Remote",
				StartDate:   "15/01/2024",
				EndDate:     "31/12/2024",
				Description: "Led development team",
			},
			wantErr: false,
		},
		{
			name: "Valid experience without optional fields",
			requestBody: dto.ExperienceRequest{
				Title:     "Developer",
				Company:   "Startup Inc",
				Type:      "Hybrid",
				StartDate: "01/06/2024",
			},
			wantErr: false,
		},
		{
			name: "Valid experience with On Site type",
			requestBody: dto.ExperienceRequest{
				Title:     "Backend Engineer",
				Company:   "Enterprise Co",
				Type:      "On Site",
				StartDate: "2023-01-01",
			},
			wantErr: false,
		},
		{
			name: "Valid experience with different date formats",
			requestBody: dto.ExperienceRequest{
				Title:     "Full Stack Engineer",
				Company:   "Tech Startup",
				Type:      "Remote",
				StartDate: "15/06/24",
				EndDate:   "31-12-24",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// Create middleware
			middleware := ValidateRequest[dto.ExperienceRequest]()

			// Execute
			middleware(c)

			// Assert
			if tt.wantErr && w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
			if !tt.wantErr && w.Code == http.StatusBadRequest {
				t.Errorf("Expected success, got status 400: %s", w.Body.String())
			}
		})
	}
}

func TestValidateRequest_InvalidExperience(t *testing.T) {
	tests := []struct {
		name        string
		requestBody map[string]interface{}
		wantStatus  int
		wantError   string
	}{
		{
			name: "Missing required title",
			requestBody: map[string]interface{}{
				"company":    "Tech Corp",
				"type":       "Remote",
				"start_date": "2024-01-01T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "title",
		},
		{
			name: "Missing required company",
			requestBody: map[string]interface{}{
				"title":      "Engineer",
				"type":       "Remote",
				"start_date": "2024-01-01T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "company",
		},
		{
			name: "Missing required type",
			requestBody: map[string]interface{}{
				"title":      "Engineer",
				"company":    "Tech Corp",
				"start_date": "2024-01-01T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "type",
		},
		{
			name: "Invalid type value",
			requestBody: map[string]interface{}{
				"title":      "Engineer",
				"company":    "Tech Corp",
				"type":       "InvalidType",
				"start_date": "2024-01-01T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "type",
		},
		{
			name: "Missing required start_date",
			requestBody: map[string]interface{}{
				"title":   "Engineer",
				"company": "Tech Corp",
				"type":    "Remote",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "start_date",
		},
		{
			name: "Invalid date format",
			requestBody: map[string]interface{}{
				"title":      "Engineer",
				"company":    "Tech Corp",
				"type":       "Remote",
				"start_date": "invalid-date",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "start_date",
		},
		{
			name: "Invalid URL format",
			requestBody: map[string]interface{}{
				"title":      "Engineer",
				"company":    "Tech Corp",
				"type":       "Remote",
				"start_date": "2024-01-01T00:00:00Z",
				"url":        "not-a-valid-url",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "url",
		},
		{
			name: "Title too long",
			requestBody: map[string]interface{}{
				"title":      string(make([]byte, 300)), // 300 characters
				"company":    "Tech Corp",
				"type":       "Remote",
				"start_date": "2024-01-01T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// Create middleware
			middleware := ValidateRequest[dto.ExperienceRequest]()

			// Execute
			middleware(c)

			// Assert
			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantError != "" {
				body := w.Body.String()
				if body == "" {
					t.Error("Expected error message in response body, got empty")
				}
			}
		})
	}
}

func TestValidateRequest_EndDateValidation(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
		wantErr   bool
	}{
		{
			name:      "End date after start date - valid",
			startDate: "01/01/2024",
			endDate:   "31/12/2024",
			wantErr:   false,
		},
		{
			name:      "End date after start date with dash format - valid",
			startDate: "01-01-2024",
			endDate:   "31-12-2024",
			wantErr:   false,
		},
		{
			name:      "End date after start date with short year - valid",
			startDate: "01/01/24",
			endDate:   "31/12/24",
			wantErr:   false,
		},
		{
			name:      "No end date - valid (current position)",
			startDate: "01/01/2024",
			endDate:   "",
			wantErr:   false,
		},
		{
			name:      "End date before start date - invalid",
			startDate: "31/12/2024",
			endDate:   "01/01/2024",
			wantErr:   true,
		},
		{
			name:      "End date same as start date - invalid",
			startDate: "01/01/2024",
			endDate:   "01/01/2024",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			requestBody := dto.ExperienceRequest{
				Title:     "Engineer",
				Company:   "Tech Corp",
				Type:      "Remote",
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}

			jsonData, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// Create middleware
			middleware := ValidateRequest[dto.ExperienceRequest]()

			// Execute
			middleware(c)

			// Assert
			if tt.wantErr && w.Code != http.StatusBadRequest {
				t.Errorf("Expected validation error (400), got %d", w.Code)
			}
			if !tt.wantErr && w.Code == http.StatusBadRequest {
				t.Errorf("Expected no error, got 400: %s", w.Body.String())
			}
		})
	}
}

func TestValidateRequest_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateRequest[dto.ExperienceRequest]()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestValidateRequest_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateRequest[dto.ExperienceRequest]()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
