package middleware

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"io"
	"net/http"

	"go.uber.org/zap"

	"crypto/sha256"
)

func CheckHashMiddleware(logger *zap.SugaredLogger, key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providedHash := r.Header.Get("HashSHA256")
			if key == "" || providedHash == "" {
				next.ServeHTTP(w, r)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Infoln("error read body")
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			providedHashBytes, err := hex.DecodeString(providedHash)
			if err != nil {
				logger.Infoln("invalid hash format")
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			h := hmac.New(sha256.New, []byte(key))
			h.Write(bodyBytes)
			expectedHash := h.Sum(nil)

			if !hmac.Equal(expectedHash, providedHashBytes) {
				http.Error(w, "", http.StatusBadRequest)
				logger.Infoln(
					"hash mismatch",
					"method", r.Method,
					"uri", r.RequestURI,
					"expectedHash", expectedHash,
					"providedHash", providedHash,
				)
				return
			}
		})
	}
}
