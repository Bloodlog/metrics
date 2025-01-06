package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	handlerLogger := logger.With("middleware", "LoggingMiddleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			handlerLogger.Infoln(
				"method", r.Method,
				"uri", r.RequestURI,
				"status", ww.Status(),
				"size", ww.BytesWritten(),
				"duration", duration,
			)
		})
	}
}
