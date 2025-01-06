package web

import (
	"metrics/internal/server/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) UpdateHandler() http.HandlerFunc {
	handlerLogger := h.logger.With("handler", "UpdateHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		var metricUpdateRequest service.MetricsUpdateRequest

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

		_, err := service.Update(metricUpdateRequest, h.memStorage)
		if err != nil {
			handlerLogger.Infow("error in service", "error", err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}
	}
}
