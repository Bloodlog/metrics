package api

import (
	"encoding/json"
	"metrics/internal/service"
	"net/http"
)

// UpdatesHandler обновляет несколько метрик.
// @Summary Обновление нескольких метрик
// @Description Обновляет несколько метрик с переданными параметрами
// @Tags Json
// @Accept  json
// @Produce  json
// @Param request body []service.MetricsUpdateRequest true "Metrics Update Request List"
// @Success 200 {string} string "Successfully updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /updates [post].
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

		err := h.metricService.UpdateMultiple(ctx, metrics)
		if err != nil {
			handlerLogger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		response.WriteHeader(http.StatusOK)
	}
}
