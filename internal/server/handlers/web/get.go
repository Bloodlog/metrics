package web

import (
	"errors"
	"metrics/internal/server/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web GetHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var metricGetRequest service.MetricsGetRequest

		metricNameRequest := chi.URLParam(request, "metricName")
		metricTypeRequest := chi.URLParam(request, "metricType")

		metricGetRequest = service.MetricsGetRequest{
			ID:    metricNameRequest,
			MType: metricTypeRequest,
		}

		metricService := service.NewMetricService(handlerLogger)
		result, err := metricService.Get(ctx, metricGetRequest, h.memStorage)
		if err != nil {
			if errors.Is(err, service.ErrMetricNotFound) {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		if metricTypeRequest == "counter" {
			_, err = response.Write([]byte(strconv.Itoa(int(*result.Delta))))
		}
		if metricTypeRequest == "gauge" {
			gaugeValue := result.Value
			result := strconv.FormatFloat(*gaugeValue, 'g', -1, 64)
			_, err = response.Write([]byte(result))
		}

		if err != nil {
			handlerLogger.Infow("error parse response", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
