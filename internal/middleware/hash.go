package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"metrics/internal/security"
	"net/http"
)

func ResponseHashMiddleware(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			hashWriter := &hashResponseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
			}
			h := hmac.New(sha256.New, []byte(key))
			h.Write(hashWriter.body.Bytes())
			hash := security.HMACSHA256Base64(hashWriter.body.Bytes(), []byte(key))

			w.Header().Set("HashSHA256", hash)

			next.ServeHTTP(hashWriter, r)
		})
	}
}

type hashResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w hashResponseWriter) Write(b []byte) (int, error) {
	n, err := w.body.Write(b)
	if err != nil {
		return n, fmt.Errorf("failed to write data to buffer: %w", err)
	}
	n, err = w.ResponseWriter.Write(b)
	if err != nil {
		return n, fmt.Errorf("failed to write response: %w", err)
	}
	return n, nil
}
