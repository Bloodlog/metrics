package api

import (
	"encoding/json"
	"errors"
	"metrics/internal/service"
	"net/http"
)

// GetHandler .
// @Summary Получение значения метрики
// @Description Получает значение метрики по имени и типу
// @Tags Json
// @Accept  json
// @Produce json
// @Param request body service.MetricsGetRequest true "Запрос на получение метрики".
// @Success 200 {object} service.MetricsResponse
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /value [post].
func (h *Handler) GetHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "api GetHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "application/json")

		var metricGetRequest service.MetricsGetRequest

		if err := json.NewDecoder(request.Body).Decode(&metricGetRequest); err != nil {
			handlerLogger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		if metricGetRequest.MType != "counter" && metricGetRequest.MType != "gauge" {
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		result, err := h.metricService.Get(ctx, metricGetRequest)
		if err != nil {
			if errors.Is(err, service.ErrMetricNotFound) {
				handlerLogger.Infoln("Metric not found", err)
				response.WriteHeader(http.StatusNotFound)
				return
			}
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
