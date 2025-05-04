package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"metrics/internal/security"
	"net/http"
	"strconv"
)

func DecryptMiddleware(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if privateKey == nil {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encryptedData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				return
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(r.Body)

			decryptedData, err := security.DecryptRSA(privateKey, encryptedData)
			if err != nil {
				http.Error(w, "failed to decrypt request body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decryptedData))
			r.ContentLength = int64(len(decryptedData))
			r.Header.Set("Content-Length", strconv.Itoa(len(decryptedData)))
			r.Header.Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		})
	}
}
