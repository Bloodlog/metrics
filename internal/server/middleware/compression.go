package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const nameError = "compression middleware"

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	n, err := g.Writer.Write(b)
	if err != nil {
		return n, fmt.Errorf("gzipResponseWriter: failed to write compressed data: %w", err)
	}
	return n, nil
}

func CompressionMiddleware(logger zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gzr, err := gzip.NewReader(r.Body)
				if err != nil {
					logger.Infoln(nameError, "Failed to decompress request body", err)
					http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
					return
				}
				defer func(gzr *gzip.Reader) {
					err := gzr.Close()
					if err != nil {
						logger.Infoln(nameError, "Failed to decompress request body", err)
					}
				}(gzr)
				r.Body = io.NopCloser(gzr)
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				logger.Infoln(nameError, "Failed to compress response", err)
				http.Error(w, "Failed to compress response", http.StatusInternalServerError)
				return
			}
			defer func(gz *gzip.Writer) {
				err := gz.Close()
				if err != nil {
					logger.Infoln(nameError, "Failed to compress response", err)
				}
			}(gz)

			w.Header().Set("Content-Encoding", "gzip")

			gzr := gzipResponseWriter{ResponseWriter: w, Writer: gz}
			next.ServeHTTP(gzr, r)
		})
	}
}
