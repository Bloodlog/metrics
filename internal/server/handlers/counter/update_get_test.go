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

func TestCounter(t *testing.T) {
	testCases := []struct {
		method       string
		path         string
		expectedBody uint64
		expectedCode int
	}{
		{
			method:       http.MethodPost,
			path:         "/update/counter/PollCount/100",
			expectedBody: uint64(100),
			expectedCode: http.StatusOK,
		},
		{
			method:       http.MethodPost,
			path:         "/update/counter/PollCount/1171",
			expectedBody: uint64(1171),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			memStorage := repository.NewMemStorage()
			r := chi.NewRouter()
			r.Post("/update/counter/{counterName}/{counterValue}", UpdateCounterHandler(memStorage))
			r.Get("/value/{metricType}/{metricName}", GetCounterHandler(memStorage))
			srv := httptest.NewServer(r)
			defer srv.Close()

			reqPost := resty.New().R()
			reqPost.Method = tc.method
			reqPost.URL = srv.URL + tc.path

			resp, err := reqPost.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Route: "+tc.method+" "+srv.URL+tc.path)

			req2 := resty.New().R()
			req2.Method = http.MethodGet
			patch := "/value/counter/PollCount"
			req2.URL = srv.URL + patch

			resp2, err := req2.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp2.StatusCode(), "Route: "+http.MethodGet+" "+srv.URL+patch)

			respBody := string(resp2.Body())
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
		})
	}
}
