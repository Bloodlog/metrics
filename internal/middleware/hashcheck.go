package middleware

import (
	"bytes"
	"io"
	"metrics/internal/security"
	"net/http"

	"go.uber.org/zap"
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

			if !security.CheckHMACSHA256Base64(bodyBytes, []byte(key), providedHash) {
				http.Error(w, "", http.StatusBadRequest)
				logger.Infoln(
					"hash mismatch",
					"method", r.Method,
					"uri", r.RequestURI,
					"providedHash", providedHash,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
