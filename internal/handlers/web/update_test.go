package web

import (
	"metrics/internal/repository"
	"metrics/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler(t *testing.T) {
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodPost, path: "/update/gauge/metricName/100", expectedCode: http.StatusOK},
		{method: http.MethodPost, path: "/update/gauge/metricName/none", expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, path: "/update/gauge/metricName", expectedCode: http.StatusNotFound},
		{method: http.MethodPost, path: "/update/counter/PollCount/100", expectedCode: http.StatusOK},
		{method: http.MethodPost, path: "/update/counter/PollCount/none", expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, path: "/update/counter/PollCount", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := logger.Sugar()

			memStorage, _ := repository.NewMemStorage()
			r := chi.NewRouter()
			metricService := service.NewMetricService(memStorage, sugar)
			webHandler := NewHandler(metricService, sugar)
			r.Post("/update/{metricType}/{metricName}/{metricValue}", webHandler.UpdateHandler())
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path
			req.SetHeader("Content-Type", "text/plain")

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
