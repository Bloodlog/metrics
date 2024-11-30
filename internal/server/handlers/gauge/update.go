package gauge

import (
	"github.com/go-chi/chi/v5"
	"metrics/internal/server/repository"
	"net/http"
	"strconv"
)

func UpdateGaugeHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		metricNameRequest := chi.URLParam(request, "metricName")
		metricValueRequest := chi.URLParam(request, "metricValue")

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		metricValue, err := strconv.ParseFloat(metricValueRequest, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		memStorage.SetGauge(metricNameRequest, metricValue)
	}
}
