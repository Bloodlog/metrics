package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestResponseCompressionMiddleware(t *testing.T) {
	logger := zap.NewExample().Sugar()
	middleware := ResponseCompressionMiddleware(logger)

	expectedBody := "compressed response data"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(expectedBody))
	})

	t.Run("With gzip encoding", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.Header.Set(acceptEncodingHeader, gzipEncoding)

		rec := httptest.NewRecorder()
		middleware(handler).ServeHTTP(rec, req)

		resp := rec.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error %v", err)
			}
		}(resp.Body)

		assert.Equal(t, gzipEncoding, resp.Header.Get(contentEncodingHeader))

		compressedBody, _ := io.ReadAll(resp.Body)

		r, _ := gzip.NewReader(bytes.NewReader(compressedBody))
		defer func(r *gzip.Reader) {
			err := r.Close()
			if err != nil {
				t.Errorf("error %v", err)
			}
		}(r)
		result, _ := io.ReadAll(r)
		decompressedBody := string(result)

		assert.Equal(t, expectedBody, decompressedBody)
	})

	t.Run("Without gzip encoding", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

		rec := httptest.NewRecorder()
		middleware(handler).ServeHTTP(rec, req)

		resp := rec.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error %v", err)
			}
		}(resp.Body)

		assert.Empty(t, resp.Header.Get(contentEncodingHeader))

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, expectedBody, string(body))
	})
}
