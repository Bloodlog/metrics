package api

import (
	"metrics/internal/server/repository"

	"go.uber.org/zap"
)

const nameLogger = "handler"
const nameError = "error"

type Handler struct {
	memStorage repository.MetricStorage
	logger     *zap.SugaredLogger
}

func NewHandler(memStorage repository.MetricStorage, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		memStorage: memStorage,
		logger:     logger,
	}
}
