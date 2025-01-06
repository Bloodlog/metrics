package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGetHandler(t *testing.T) {
	counterValue := uint64(100)
	gaugeValue := 1234.1234

	testCases := []struct {
		name         string
		requestBody  string
		setupStorage func(memStorage *repository.MemStorage)
		expectedBody string
		expectedCode int
	}{
		{
			name:        "Get Counter Successfully",
			requestBody: `{"id": "PollCount", "type": "counter"}`,
			setupStorage: func(memStorage *repository.MemStorage) {
				err := memStorage.SetCounter("PollCount", counterValue)
				if err != nil {
					t.Errorf("Failed to SetCounter: %v", err)
					return
				}
			},
			expectedBody: `{"id":"PollCount","type":"counter","delta":100}`,
			expectedCode: http.StatusOK,
		},
		{
			name:        "Get Gauge Successfully",
			requestBody: `{"id": "Allocate", "type": "gauge"}`,
			setupStorage: func(memStorage *repository.MemStorage) {
				err := memStorage.SetGauge("Allocate", gaugeValue)
				if err != nil {
					t.Errorf("Failed to SetCounter: %v", err)
					return
				}
			},
			expectedBody: `{"id":"Allocate","type":"gauge","value":1234.1234}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Metric Type",
			requestBody:  `{"id": "Unknown", "type": "invalid"}`,
			setupStorage: func(memStorage *repository.MemStorage) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Metric Not Found",
			requestBody:  `{"id": "Unknown", "type": "counter"}`,
			setupStorage: func(memStorage *repository.MemStorage) {},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid JSON",
			requestBody:  `{"id": nil, "type": "counter", "delta": "invalid"}`,
			setupStorage: func(memStorage *repository.MemStorage) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := logger.Sugar()
			memStorage := repository.NewMemStorage()
			tc.setupStorage(memStorage)

			r := chi.NewRouter()
			apiHandler := NewHandler(memStorage, sugar)
			r.Post("/get", apiHandler.GetHandler())
			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.requestBody).
				Post(srv.URL + "/get")

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Unexpected status code")

			if tc.expectedCode == http.StatusOK {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()), "Unexpected response body")
			}
		})
	}
}
