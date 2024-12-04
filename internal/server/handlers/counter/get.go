package counter

import (
	"fmt"
	"metrics/internal/server/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetCounterHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		metricNameRequest := chi.URLParam(request, "metricName")

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		counter, err := memStorage.GetCounter(metricNameRequest)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		_, writeErr := response.Write([]byte(fmt.Sprintf("%d", counter)))
		if writeErr != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
