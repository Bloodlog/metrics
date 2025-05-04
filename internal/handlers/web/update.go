package web

import (
	"metrics/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// UpdateHandler .
// @Summary Обновление значения метрики
// @Description Обновляет значение метрики по её типу (counter или gauge) на основе переданных параметров
// @Tags Text
// @Accept  text/plain
// @Produce  text/plain
// @Param metricType path string true "Тип метрики (counter или gauge)"
// @Param metricName path string true "Имя метрики"
// @Param metricValue path string true "Новое значение метрики"
// @Success 200 {string} string "Метрика успешно обновлена"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /update/{metricType}/{metricName}/{metricValue} [post].
func (h *Handler) UpdateHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web UpdateHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		var metricUpdateRequest service.MetricsUpdateRequest

		metricValueRequest := chi.URLParam(request, "metricValue")
		metricNameRequest := chi.URLParam(request, "metricName")
		metricTypeRequest := chi.URLParam(request, "metricType")

		if metricTypeRequest == "counter" {
			metricValue, err := strconv.ParseInt(metricValueRequest, 10, 64)
			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				return
			}

			metricUpdateRequest = service.MetricsUpdateRequest{
				ID:    metricNameRequest,
				MType: metricTypeRequest,
				Delta: &metricValue,
			}
		}
		if metricTypeRequest == "gauge" {
			metricValue, err := strconv.ParseFloat(metricValueRequest, 64)
			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				return
			}
			metricUpdateRequest = service.MetricsUpdateRequest{
				ID:    metricNameRequest,
				MType: metricTypeRequest,
				Value: &metricValue,
			}
		}
		_, err := h.metricService.Update(ctx, metricUpdateRequest)
		if err != nil {
			handlerLogger.Infow("error in service", "error", err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}
	}
}
