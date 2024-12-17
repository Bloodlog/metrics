package web

import (
	"errors"
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func GetHandler(memStorage *repository.MemStorage, logger zap.SugaredLogger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		const nameError = "error"
		var metricGetRequest service.MetricsGetRequest

		metricNameRequest := chi.URLParam(request, "metricName")
		metricTypeRequest := chi.URLParam(request, "metricType")

		metricGetRequest = service.MetricsGetRequest{
			ID:    metricNameRequest,
			MType: metricTypeRequest,
		}

		result, err := service.Get(metricGetRequest, memStorage)
		if err != nil {
			if errors.Is(err, service.ErrMetricNotFound) {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Infow("error in service", nameError, err)
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
			logger.Infow("error parse response", "error", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
