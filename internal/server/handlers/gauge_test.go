package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGaugeHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	srv := httptest.NewServer(GaugeHandler(memStorage, false))
	defer srv.Close()

	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodPost, path: "/update/gauge/metricName/100", expectedCode: http.StatusOK},
		{method: http.MethodPost, path: "/update/gauge/metricName/none", expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, path: "/update/gauge/metricName", expectedCode: http.StatusNotFound},
		{method: http.MethodGet, path: "/update/uknown/metricName/100", expectedCode: http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
