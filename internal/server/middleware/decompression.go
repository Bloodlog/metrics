package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func DecompressionMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	handlerLogger := logger.With("middleware", "DecompressionMiddleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get(contentEncodingHeader), gzipEncoding) {
				gzr, err := gzip.NewReader(r.Body)
				if err != nil {
					handlerLogger.Infoln("Failed to decompress request body", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer func(gzr *gzip.Reader) {
					err := gzr.Close()
					if err != nil {
						handlerLogger.Infoln("Failed to decompress request body", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}(gzr)
				r.Body = io.NopCloser(gzr)
			}

			next.ServeHTTP(w, r)
		})
	}
}
