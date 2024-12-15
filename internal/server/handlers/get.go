package handlers

import (
	"encoding/json"
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"

	"go.uber.org/zap"
)

func GetHandler(memStorage *repository.MemStorage, logger zap.SugaredLogger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		const nameError = "error"

		var metricGetRequest service.MetricsGetRequest
		if err := json.NewDecoder(request.Body).Decode(&metricGetRequest); err != nil {
			logger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		result, err := service.Get(metricGetRequest, memStorage)
		if err != nil {
			logger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		resp, err := json.Marshal(result)
		if err != nil {
			logger.Infow("error marshal json", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = response.Write(resp)
		if err != nil {
			logger.Infow("error write response", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
