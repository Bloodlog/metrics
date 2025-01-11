package api

import (
	"encoding/json"
	"metrics/internal/server/service"
	"net/http"
)

func (h *Handler) UpdateHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "api UpdateHandler")
	const nameError = "error"
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "application/json")

		var metricUpdateRequest service.MetricsUpdateRequest

		if err := json.NewDecoder(request.Body).Decode(&metricUpdateRequest); err != nil {
			handlerLogger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		metricService := service.NewMetricService(handlerLogger)
		result, err := metricService.Update(ctx, metricUpdateRequest, h.memStorage)
		if err != nil {
			handlerLogger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(result)
		if err != nil {
			handlerLogger.Infow("error marshal json", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = response.Write(resp)
		if err != nil {
			handlerLogger.Infow("error write response", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
