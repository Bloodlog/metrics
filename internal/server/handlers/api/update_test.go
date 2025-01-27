package api

import (
	"context"
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
			ctx := context.Background()

			memStorage, _ := repository.NewMemStorage(ctx)
			r := chi.NewRouter()
			apiHandler := NewHandler(memStorage, sugar)
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
