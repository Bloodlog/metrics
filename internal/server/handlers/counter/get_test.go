package counter

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGetCounterHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	counterValue := uint64(100)
	memStorage.SetCounter("PollCount", counterValue)
	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", GetCounterHandler(memStorage))
	srv := httptest.NewServer(r)
	defer srv.Close()

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
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			respBody := string(resp.Body())
			if respBody == "" {
				t.Fatalf("тело ответа пустое")
			}

			bodyUint64, err := strconv.ParseUint(respBody, 10, 64)
			if err != nil {
				t.Fatalf("не удалось преобразовать тело в uint64: %v", err)
			}
			assert.Equal(t, tc.expectedBody, bodyUint64)
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}

func TestGetCounterFailsHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	counterValue := uint64(100)
	memStorage.SetCounter("PollCount", counterValue)
	r := chi.NewRouter()
	r.Get("/value/counter/{metricName}", GetCounterHandler(memStorage))
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

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
