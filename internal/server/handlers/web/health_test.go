package web

import (
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthHandler_Failure(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	memStorage, _ := repository.NewMemStorage()
	metricService := service.NewMetricService(memStorage, sugar)
	webHandler := NewHandler(metricService, sugar)
	healthHandler := webHandler.HealthHandler("mock_dsn")

	req := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
	w := httptest.NewRecorder()

	healthHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
