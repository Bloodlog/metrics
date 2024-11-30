package gauge

import (
	"github.com/go-chi/chi/v5"
	"metrics/internal/server/repository"
	"net/http"
	"strconv"
)

func GetGaugeHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		metricNameRequest := chi.URLParam(request, "metricName")

		gauge, err := memStorage.GetGauge(metricNameRequest)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		_, writeErr := response.Write([]byte(strconv.FormatFloat(gauge, 'g', -1, 64)))
		if writeErr != nil {
			http.Error(response, "failed to write response", http.StatusInternalServerError)
		}

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)
	}
}
