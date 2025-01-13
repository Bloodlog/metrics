package api

import (
	"encoding/json"
	"metrics/internal/server/service"
	"net/http"
)

func (h *Handler) UpdatesHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "api UpdateHandler")
	const nameError = "error"
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "application/json")

		var metrics []service.MetricsUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&metrics); err != nil {
			handlerLogger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		metricService := service.NewMetricService(handlerLogger)

		err := metricService.UpdateMultiple(ctx, metrics, h.memStorage)
		if err != nil {
			handlerLogger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		response.WriteHeader(http.StatusOK)
	}
}
