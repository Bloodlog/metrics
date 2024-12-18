package api

import (
	"encoding/json"
	"metrics/internal/server/repository"
	"metrics/internal/server/service"
	"net/http"

	"go.uber.org/zap"
)

func UpdateHandler(memStorage repository.MetricStorage, logger zap.SugaredLogger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		const nameError = "update handler"

		var metricUpdateRequest service.MetricsUpdateRequest

		if err := json.NewDecoder(request.Body).Decode(&metricUpdateRequest); err != nil {
			logger.Infow("Invalid JSON", "error", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := service.Update(metricUpdateRequest, memStorage)
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
