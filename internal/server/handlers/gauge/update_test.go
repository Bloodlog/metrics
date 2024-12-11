package gauge

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestGaugeHandler(t *testing.T) {
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodPost, path: "/update/gauge/metricName/100", expectedCode: http.StatusOK},
		{method: http.MethodPost, path: "/update/gauge/metricName/none", expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, path: "/update/gauge/metricName", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			memStorage := repository.NewMemStorage()
			r := chi.NewRouter()
			r.Post("/update/gauge/{metricName}/{metricValue}", UpdateGaugeHandler(memStorage))
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
