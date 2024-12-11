package counter

import (
	"log"
	"metrics/internal/server/repository"
	"net/http"
	"strconv"

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

		counterString := strconv.Itoa(int(counter))
		_, err = response.Write([]byte(counterString))
		if err != nil {
			log.Printf("error get counter: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
