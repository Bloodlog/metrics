package counter

import (
	"metrics/internal/server/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func UpdateCounterHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		counterNameRequest := chi.URLParam(request, "counterName")
		counterValueRequest := chi.URLParam(request, "counterValue")

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		counterValue, err := strconv.ParseUint(counterValueRequest, 10, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		memStorage.SetCounter(counterNameRequest, counterValue)
	}
}
