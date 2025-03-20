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

func gzipCompress(data string) *bytes.Buffer {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte(data))
	_ = gz.Close()
	return &buf
}

func TestDecompressionMiddleware(t *testing.T) {
	logger := zap.NewExample().Sugar()
	middleware := DecompressionMiddleware(logger)

	expectedBody := "test data"

	compressedBody := gzipCompress(expectedBody)

	req := httptest.NewRequest(http.MethodPost, "/", compressedBody)
	req.Header.Set("Content-Encoding", "gzip")

	var actualBody string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		actualBody = string(body)
	})

	middleware(handler).ServeHTTP(httptest.NewRecorder(), req)

	assert.Equal(t, expectedBody, actualBody)
}
