package counter

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"metrics/internal/server/repository"
	"net/http"
)

func GetCounterHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		metricNameRequest := chi.URLParam(request, "metricName")

		counter, err := memStorage.GetCounter(metricNameRequest)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		_, writeErr := response.Write([]byte(fmt.Sprintf("%d", counter)))
		if writeErr != nil {
			http.Error(response, "failed to write response", http.StatusInternalServerError)
		}

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)
	}
}
