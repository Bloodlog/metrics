package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestBatchUpdateHandler(t *testing.T) {
	successCases := []struct {
		name         string
		requestBody  string
		expectedCode int
		metricID     string
	}{
		{
			name: "Valid Counter and Gauge Update",
			requestBody: `[
				{"id": "PollCount", "type": "counter", "delta": 100},
				{"id": "Allocate", "type": "gauge", "value": 25.5}
			]`,
			expectedCode: http.StatusOK,
		},
		{
			name: "Valid Counter and Gauge Update (Allocate)",
			requestBody: `[
				{"id": "PollCount", "type": "counter", "delta": 100},
				{"id": "Allocate", "type": "gauge", "value": 25.5}
			]`,
			expectedCode: http.StatusOK,
			metricID:     "Allocate",
		},
	}

	failCases := []struct {
		name         string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "Invalid JSON",
			requestBody:  `{"id": "PollCount", "type": "counter", "delta": "not_a_number"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing Required Fields",
			requestBody:  `{"id": "PollCount"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Unsupported Metric Type",
			requestBody:  `{"id": "Unknown", "type": "unknown"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range successCases {
		t.Run(tc.name, func(t *testing.T) {
			route := "/updates"
			r, apiHandler, memStore := setupHandler()
			r.Post(route, apiHandler.UpdatesHandler())
			srv := httptest.NewServer(r)
			defer srv.Close()
			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + route)

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")
			ctx := context.Background()
			if tc.expectedCode == http.StatusOK {
				if tc.metricID == "PollCount" {
					counter, err := memStore.GetCounter(ctx, tc.metricID)
					assert.NoError(t, err, "Error retrieving counter")
					assert.Equal(t, 100, counter, "Unexpected counter value")
				}
				if tc.metricID == "Allocate" {
					gauge, err := memStore.GetGauge(ctx, tc.metricID)
					assert.NoError(t, err, "Error retrieving gauge")
					assert.Equal(t, 25.5, gauge, "Unexpected gauge value")
				}
			}
		})
	}

	for _, tc := range failCases {
		t.Run(tc.name, func(t *testing.T) {
			r, apiHandler, _ := setupHandler()
			r.Post("/updates", apiHandler.UpdatesHandler())

			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + "/updates")

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")
		})
	}
}
