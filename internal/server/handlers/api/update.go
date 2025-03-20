package api

import (
	"encoding/json"
	"metrics/internal/server/dto"
	"net/http"
)

// UpdateHandler .
// @Summary Обновление метрики
// @Description Обновляет метрику с переданными параметрами
// @Tags Json
// @Accept  json
// @Produce  json
// @Param request body dto.MetricsUpdateRequest true "Metrics Update Request"
// @Success 200 {object} string "Response with success status"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /update [post].
func (h *Handler) UpdateHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "api UpdateHandler")
	const nameError = "error"
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "application/json")

		var metricUpdateRequest dto.MetricsUpdateRequest

		if err := json.NewDecoder(request.Body).Decode(&metricUpdateRequest); err != nil {
			handlerLogger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := h.metricService.Update(ctx, metricUpdateRequest)
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
