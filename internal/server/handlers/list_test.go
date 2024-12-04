package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestListGaugeHandler(t *testing.T) {
	memStorage := repository.NewMemStorage()
	metricName := "metricName"
	metricValue := 1234.1234
	memStorage.SetGauge(metricName, metricValue)
	counterValue := uint64(100)
	counterName := "PollCount"
	memStorage.SetCounter(counterName, counterValue)
	r := chi.NewRouter()
	r.Get("/", ListHandler(memStorage))
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		method       string
		path         string
		expectedBody float64
		expectedCode int
	}{
		{method: http.MethodGet, path: "/", expectedBody: metricValue, expectedCode: http.StatusOK},
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
				t.Fatalf("response body is empty")
			}

			metricNameMatch, _ := regexp.MatchString(metricName, respBody)

			metricValueStr := fmt.Sprintf("%f", metricValue)
			numberMatch, _ := regexp.MatchString(metricValueStr, respBody)
			assert.True(t, metricNameMatch, "metric Name name is not exist on page")
			assert.True(t, numberMatch, "metric Value is not exist on page")

			counterNameMatch, _ := regexp.MatchString(counterName, respBody)

			counterValueStr := fmt.Sprintf("%d", counterValue)
			counterMatch, _ := regexp.MatchString(counterValueStr, respBody)

			assert.True(t, counterNameMatch, "counter name is not exist on page")
			assert.True(t, counterMatch, "counter value is not exist on page")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected. Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
