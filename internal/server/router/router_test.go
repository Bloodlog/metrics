package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"metrics/internal/server/repository"
)

func TestRouter(t *testing.T) {
	memStorage := repository.NewMemStorage()
	router := chi.NewRouter()
	register(router, memStorage, false)
	srv := httptest.NewServer(router)
	defer srv.Close()

	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodGet, path: "/", expectedCode: http.StatusOK},
		{method: http.MethodPost, path: "/update/unknown/testCounter/100", expectedCode: http.StatusBadRequest},
		{method: http.MethodGet, path: "/value/unknown/testCounter", expectedCode: http.StatusBadRequest},
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
