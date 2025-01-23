package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func ResponseHashMiddleware(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			hashWriter := &hashResponseWriter{ResponseWriter: w}
			next.ServeHTTP(hashWriter, r)

			if len(hashWriter.body) > 0 {
				h := hmac.New(sha256.New, []byte(key))
				h.Write(hashWriter.body)
				hashHex := hex.EncodeToString(h.Sum(nil))

				w.Header().Set("HashSHA256", hashHex)
			}
		})
	}
}

type hashResponseWriter struct {
	http.ResponseWriter
	body []byte
}

func (w *hashResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)

	return w.ResponseWriter.Write(data)
}

func (w *hashResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *hashResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
