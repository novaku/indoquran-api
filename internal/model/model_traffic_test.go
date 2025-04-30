package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestApiTrafficLog(t *testing.T) {
	tests := []struct {
		name     string
		traffic  ApiTrafficLog
		expected ApiTrafficLog
	}{
		{
			name: "Valid API traffic log entry",
			traffic: ApiTrafficLog{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				IPAddress:      "192.168.1.1",
				Endpoint:       "/api/v1/test",
				Duration:       "150ms",
				HTTPMethod:     "GET",
				RequestPayload: `{"key": "value"}`,
				ResponseStatus: 200,
				ResponseBody:   `{"status": "success"}`,
				UserAgent:      "Mozilla/5.0",
				Referer:        "https://example.com",
			},
			expected: ApiTrafficLog{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				IPAddress:      "192.168.1.1",
				Endpoint:       "/api/v1/test",
				Duration:       "150ms",
				HTTPMethod:     "GET",
				RequestPayload: `{"key": "value"}`,
				ResponseStatus: 200,
				ResponseBody:   `{"status": "success"}`,
				UserAgent:      "Mozilla/5.0",
				Referer:        "https://example.com",
			},
		},
		{
			name: "With deleted at",
			traffic: ApiTrafficLog{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: gorm.DeletedAt{Time: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
				},
				IPAddress:      "10.0.0.1",
				Endpoint:       "/api/v1/users",
				Duration:       "200ms",
				HTTPMethod:     "POST",
				RequestPayload: `{"username": "test"}`,
				ResponseStatus: 201,
				ResponseBody:   `{"id": 1}`,
				UserAgent:      "PostmanRuntime/7.28.4",
				Referer:        "",
			},
			expected: ApiTrafficLog{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: gorm.DeletedAt{Time: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
				},
				IPAddress:      "10.0.0.1",
				Endpoint:       "/api/v1/users",
				Duration:       "200ms",
				HTTPMethod:     "POST",
				RequestPayload: `{"username": "test"}`,
				ResponseStatus: 201,
				ResponseBody:   `{"id": 1}`,
				UserAgent:      "PostmanRuntime/7.28.4",
				Referer:        "",
			},
		},
		{
			name: "Error response log",
			traffic: ApiTrafficLog{
				IPAddress:      "127.0.0.1",
				Endpoint:       "/api/v1/invalid",
				Duration:       "50ms",
				HTTPMethod:     "DELETE",
				RequestPayload: "",
				ResponseStatus: 404,
				ResponseBody:   `{"error": "Not Found"}`,
				UserAgent:      "curl/7.64.1",
				Referer:        "",
			},
			expected: ApiTrafficLog{
				IPAddress:      "127.0.0.1",
				Endpoint:       "/api/v1/invalid",
				Duration:       "50ms",
				HTTPMethod:     "DELETE",
				RequestPayload: "",
				ResponseStatus: 404,
				ResponseBody:   `{"error": "Not Found"}`,
				UserAgent:      "curl/7.64.1",
				Referer:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.ID, tt.traffic.ID)
			assert.Equal(t, tt.expected.CreatedAt, tt.traffic.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, tt.traffic.UpdatedAt)
			assert.Equal(t, tt.expected.DeletedAt, tt.traffic.DeletedAt)
			assert.Equal(t, tt.expected.IPAddress, tt.traffic.IPAddress)
			assert.Equal(t, tt.expected.Endpoint, tt.traffic.Endpoint)
			assert.Equal(t, tt.expected.Duration, tt.traffic.Duration)
			assert.Equal(t, tt.expected.HTTPMethod, tt.traffic.HTTPMethod)
			assert.Equal(t, tt.expected.RequestPayload, tt.traffic.RequestPayload)
			assert.Equal(t, tt.expected.ResponseStatus, tt.traffic.ResponseStatus)
			assert.Equal(t, tt.expected.ResponseBody, tt.traffic.ResponseBody)
			assert.Equal(t, tt.expected.UserAgent, tt.traffic.UserAgent)
			assert.Equal(t, tt.expected.Referer, tt.traffic.Referer)
		})
	}
}
