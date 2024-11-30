package gauge

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGetGaugeHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	metricValue := 1234.1234
	memStorage.SetGauge("metricName", metricValue)
	r := chi.NewRouter()
	r.Get("/value/gauge/{metricName}", GetGaugeHandler(memStorage))
	srv := httptest.NewServer(r)
	defer srv.Close()

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
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			respBody := string(resp.Body())
			if respBody == "" {
				t.Fatalf("тело ответа пустое")
			}

			assert.Equal(t, strconv.FormatFloat(tc.expectedBody, 'g', -1, 64), respBody)
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected. Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}

func TestGetGaugeFailHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	metricValue := 1234.1234
	memStorage.SetGauge("metricName", metricValue)
	r := chi.NewRouter()
	r.Get("/value/gauge/{metricName}", GetGaugeHandler(memStorage))
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodGet, path: "/value/gauge/unknown", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected. Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
