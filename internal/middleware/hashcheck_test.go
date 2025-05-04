package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCheckHashMiddleware(t *testing.T) {
	logger := zap.NewNop().Sugar()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	key := "secret"
	body := `{"key":"value"}`
	hash := generateHMACSHA256Hash(body, key)
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	tests := []struct {
		name               string
		providedHash       string
		expectedStatusCode int
	}{
		{
			name:               "Valid hash",
			providedHash:       encodedHash,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "No hash provided",
			providedHash:       "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid hash format",
			providedHash:       "invalid_hash",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Hash mismatch",
			providedHash:       base64.StdEncoding.EncodeToString([]byte("wronghash")),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
			req.Header.Set("HashSHA256", tt.providedHash)

			rr := httptest.NewRecorder()

			middleware := CheckHashMiddleware(logger, key)

			middleware(nextHandler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func generateHMACSHA256Hash(data, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}
