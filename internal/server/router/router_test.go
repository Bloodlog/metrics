package router

import (
	"context"
	"metrics/internal/server/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"metrics/internal/server/repository"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{method: http.MethodGet, path: "/", expectedCode: http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			configs := config.Config{}
			logger := zap.NewNop()
			sugar := logger.Sugar()
			ctx := context.Background()

			memStorage, _ := repository.NewMemStorage(ctx)
			router := chi.NewRouter()
			register(router, &configs, memStorage, sugar)
			srv := httptest.NewServer(router)
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
