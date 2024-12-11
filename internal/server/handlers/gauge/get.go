package gauge

import (
	"log"
	"metrics/internal/server/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetGaugeHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		metricNameRequest := chi.URLParam(request, "metricName")
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		gauge, err := memStorage.GetGauge(metricNameRequest)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = response.Write([]byte(strconv.FormatFloat(gauge, 'g', -1, 64)))
		if err != nil {
			log.Printf("error get metric: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
