package web

import (
	"metrics/internal/server/repository"

	"go.uber.org/zap"
)

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
