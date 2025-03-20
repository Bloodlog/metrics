package api

import (
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
	testCases := []struct {
		name         string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "Valid Counter Update",
			requestBody:  `{"id": "PollCount", "type": "counter", "delta": 100}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid Gauge Update",
			requestBody:  `{"id": "Allocate", "type": "gauge", "value": 25.5}`,
			expectedCode: http.StatusOK,
		},
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := logger.Sugar()

			memStorage, _ := repository.NewMemStorage()
			r := chi.NewRouter()
			metricService := service.NewMetricService(memStorage, sugar)
			apiHandler := NewHandler(metricService, sugar)
			r.Post("/update", apiHandler.UpdateHandler())

			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + "/update")

			assert.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")
		})
	}
}
