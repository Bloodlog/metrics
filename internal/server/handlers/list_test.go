package handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"metrics/internal/server/repository"
)

func TestListGaugeHandler(t *testing.T) {
	metricValue := 1234.1234

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
			memStorage := repository.NewMemStorage()
			metricName := "metricName"
			memStorage.SetGauge(metricName, metricValue)
			counterValue := uint64(100)
			counterName := "PollCount"
			memStorage.SetCounter(counterName, counterValue)
			r := chi.NewRouter()
			r.Get("/", ListHandler(memStorage))
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			respBody := string(resp.Body())
			if respBody == "" {
				t.Error("response body is empty")
				return
			}

			metricValueStr := strconv.FormatFloat(metricValue, 'f', -1, 64)

			assert.Contains(t, respBody, metricName, "metric Name is not exist on page")
			assert.Contains(t, respBody, metricValueStr, "metric Value is not exist on page")

			counterNameMatch, _ := regexp.MatchString(counterName, respBody)

			counterValueStr := strconv.FormatUint(counterValue, 10)
			counterMatch, err := regexp.MatchString(counterValueStr, respBody)
			if err != nil {
				t.Error("parsing regexp in response body")
				return
			}

			assert.True(t, counterNameMatch, "counter name is not exist on page")
			assert.True(t, counterMatch, "counter value is not exist on page")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)
		})
	}
}
