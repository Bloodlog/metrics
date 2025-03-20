package web

import (
	"errors"
	"metrics/internal/server/apperrors"
	"metrics/internal/server/dto"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetHandler .
// @Summary Получение значения метрики
// @Description Возвращает значение метрики по ее типу (counter или gauge) в формате текста
// @Tags Text
// @Accept  text/plain
// @Produce  text/plain
// @Param metricType path string true "Тип метрики (counter или gauge)"
// @Param metricName path string true "Имя метрики"
// @Success 200 {string} string "Метрика возвращена успешно"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /value/{metricType}/{metricName} [get].
func (h *Handler) GetHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web GetHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var metricGetRequest dto.MetricsGetRequest

		metricNameRequest := chi.URLParam(request, "metricName")
		metricTypeRequest := chi.URLParam(request, "metricType")

		metricGetRequest = dto.MetricsGetRequest{
			ID:    metricNameRequest,
			MType: metricTypeRequest,
		}

		result, err := h.metricService.Get(ctx, metricGetRequest)
		if err != nil {
			if errors.Is(err, apperrors.ErrMetricNotFound) {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		if metricTypeRequest == "counter" {
			_, err = response.Write([]byte(strconv.Itoa(int(*result.Delta))))
		}
		if metricTypeRequest == "gauge" {
			gaugeValue := result.Value
			result := strconv.FormatFloat(*gaugeValue, 'g', -1, 64)
			_, err = response.Write([]byte(result))
		}

		if err != nil {
			handlerLogger.Infow("error parse response", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
