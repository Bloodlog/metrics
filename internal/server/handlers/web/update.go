package web

import (
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func UpdateHandler(memStorage repository.MetricStorage, logger zap.SugaredLogger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		var metricUpdateRequest service.MetricsUpdateRequest
		const nameError = "error"

		metricValueRequest := chi.URLParam(request, "metricValue")
		metricNameRequest := chi.URLParam(request, "metricName")
		metricTypeRequest := chi.URLParam(request, "metricType")

		if metricTypeRequest == "counter" {
			metricValue, err := strconv.ParseInt(metricValueRequest, 10, 64)
			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				return
			}

			metricUpdateRequest = service.MetricsUpdateRequest{
				ID:    metricNameRequest,
				MType: metricTypeRequest,
				Delta: &metricValue,
			}
		}
		if metricTypeRequest == "gauge" {
			metricValue, err := strconv.ParseFloat(metricValueRequest, 64)
			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				return
			}
			metricUpdateRequest = service.MetricsUpdateRequest{
				ID:    metricNameRequest,
				MType: metricTypeRequest,
				Value: &metricValue,
			}
		}

		_, err := service.Update(metricUpdateRequest, memStorage)
		if err != nil {
			logger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}
	}
}
