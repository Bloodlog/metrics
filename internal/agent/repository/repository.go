package repository

import "metrics/internal/agent/dto"

type MetricsRepository interface {
	GetMetrics() []dto.Metric
}
