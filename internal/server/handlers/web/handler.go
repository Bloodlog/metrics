package web

import (
	"metrics/internal/server/service"

	"go.uber.org/zap"
)

const nameLogger = "handler"
const nameError = "error"

type Handler struct {
	metricService service.MetricService
	logger        *zap.SugaredLogger
}

func NewHandler(
	metricService service.MetricService,
	logger *zap.SugaredLogger,
) *Handler {
	return &Handler{
		metricService: metricService,
		logger:        logger,
	}
}
