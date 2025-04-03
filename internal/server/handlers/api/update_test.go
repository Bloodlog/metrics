package api

import (
	"context"
	"metrics/internal/server/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestUpdateHandler(t *testing.T) {
	successCases := []struct {
		name         string
		requestBody  string
		expectedCode int
		metricID     string
	}{
		{
			name:         "Valid Counter Update",
			requestBody:  `{"id": "PollCount", "type": "counter", "delta": 100}`,
			expectedCode: http.StatusOK,
			metricID:     "PollCount",
		},
		{
			name:         "Valid Gauge Update",
			requestBody:  `{"id": "Allocate", "type": "gauge", "value": 25.5}`,
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
			route := "/update"
			r, apiHandler, memStore := setupHandler()
			r.Post(route, apiHandler.UpdateHandler())
			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + route)

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")

			if tc.expectedCode == http.StatusOK {
				ctx := context.Background()
				if tc.requestBody[10:15] == "gauge" {
					gauge, err := memStore.GetGauge(ctx, tc.metricID)
					assert.NoError(t, err, "Error retrieving gauge")
					assert.Equal(t, 25.5, gauge, "Unexpected gauge value")
				}
				if tc.requestBody[10:17] == "counter" {
					counter, err := memStore.GetCounter(ctx, tc.metricID)
					assert.NoError(t, err, "Error retrieving counter")
					assert.Equal(t, 100, counter, "Unexpected counter value")
				}
			}
		})
	}

	for _, tc := range failCases {
		t.Run(tc.name, func(t *testing.T) {
			route := "/update"
			r, apiHandler, _ := setupHandler()
			r.Post(route, apiHandler.UpdateHandler())

			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + route)

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")
		})
	}
}

func setupHandler() (*chi.Mux, *Handler, repository.MetricStorage) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	memStorage, _ := repository.NewMemStorage()
	r := chi.NewRouter()
	metricService := service.NewMetricService(memStorage, sugar)
	apiHandler := NewHandler(metricService, sugar)

	return r, apiHandler, memStorage
}
