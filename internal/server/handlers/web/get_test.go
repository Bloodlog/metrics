package web

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGetCounterHandler(t *testing.T) {
	counterValue := uint64(100)

	testCases := []struct {
		method       string
		path         string
		expectedBody uint64
		expectedCode int
	}{
		{method: http.MethodGet, path: "/value/counter/PollCount", expectedBody: counterValue, expectedCode: http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := *logger.Sugar()

			memStorage := repository.NewMemStorage()
			memStorage.SetCounter("PollCount", counterValue)
			r := chi.NewRouter()
			r.Get("/value/{metricType}/{metricName}", GetHandler(memStorage, sugar))
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path
			req.SetHeader("Content-Type", "text/plain")

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			respBody := string(resp.Body())
			if respBody == "" {
				t.Error("response body is empty")
				return
			}

			bodyUint64, err := strconv.ParseUint(respBody, 10, 64)
			if err != nil {
				t.Error("не удалось преобразовать тело в uint64")
				return
			}
			assert.Equal(t, tc.expectedBody, bodyUint64)
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}

func TestGetCounterFailsHandler(t *testing.T) {
	logger := zap.NewNop()
	sugar := *logger.Sugar()

	memStorage := repository.NewMemStorage()
	counterValue := uint64(100)
	memStorage.SetCounter("PollCount", counterValue)
	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", GetHandler(memStorage, sugar))
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodGet, path: "/value/counter/unknown", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
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

func TestGetGaugeHandler(t *testing.T) {
	metricValue := 1234.1234
	testCases := []struct {
		method       string
		path         string
		expectedBody float64
		expectedCode int
	}{
		{method: http.MethodGet, path: "/value/gauge/metricName", expectedBody: metricValue, expectedCode: http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := *logger.Sugar()

			memStorage := repository.NewMemStorage()
			memStorage.SetGauge("metricName", metricValue)
			r := chi.NewRouter()
			r.Get("/value/{metricType}/{metricName}", GetHandler(memStorage, sugar))
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path
			req.SetHeader("Content-Type", "text/plain")

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			respBody := string(resp.Body())
			if respBody == "" {
				t.Error("response body is empty")
				return
			}

			assert.Equal(t, strconv.FormatFloat(tc.expectedBody, 'g', -1, 64), respBody)
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}

func TestGetGaugeFailHandler(t *testing.T) {
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodGet, path: "/value/gauge/unknown", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			logger := zap.NewNop()
			sugar := *logger.Sugar()

			memStorage := repository.NewMemStorage()
			metricValue := 1234.1234
			memStorage.SetGauge("metricName", metricValue)
			r := chi.NewRouter()
			r.Get("/value/{metricType}/{metricName}", GetHandler(memStorage, sugar))
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