package web

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

// Health godoc
// @Tags Info
// @Summary Запрос состояния сервиса
// @ID infoHealth
// @Accept  json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /ping [get]
func (h *Handler) HealthHandler(dsn string) http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web HealthHandler")
	const nameError = "error"
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		_, err := pgx.Connect(ctx, dsn)
		if err != nil {
			handlerLogger.Infow("Unable to connect to database", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
