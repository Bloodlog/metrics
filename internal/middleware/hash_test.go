package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseHashMiddleware(t *testing.T) {
	const secretKey = "secret"
	expectedBody := "test response"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(expectedBody))
	})

	t.Run("With secret key", func(t *testing.T) {
		middleware := ResponseHashMiddleware(secretKey)
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		resp := rec.Result()
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("error %v", err)
			}
		}()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, expectedBody, string(body))

		hash := resp.Header.Get("HashSHA256")
		assert.NotEmpty(t, hash)
	})

	t.Run("Without secret key", func(t *testing.T) {
		middleware := ResponseHashMiddleware("")
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		resp := rec.Result()
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("error %v", err)
			}
		}()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, expectedBody, string(body))

		hash := resp.Header.Get("HashSHA256")
		assert.Empty(t, hash)
	})
}
