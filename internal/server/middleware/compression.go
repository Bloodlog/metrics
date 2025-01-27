package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

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

const (
	contentEncodingHeader = "Content-Encoding"
	acceptEncodingHeader  = "Accept-Encoding"
	gzipEncoding          = "gzip"
)

func CompressionMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	handlerLogger := logger.With("middleware", "CompressionMiddleware")
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

			if !strings.Contains(r.Header.Get(acceptEncodingHeader), gzipEncoding) {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				handlerLogger.Infoln("Failed to compress response", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer func(gz *gzip.Writer) {
				err := gz.Close()
				if err != nil {
					handlerLogger.Infoln("Failed to compress response", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}(gz)

			w.Header().Set(contentEncodingHeader, gzipEncoding)

			gzr := gzipResponseWriter{ResponseWriter: w, Writer: gz}
			next.ServeHTTP(gzr, r)
		})
	}
}
