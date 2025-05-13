package rpc

import (
	"context"
	"fmt"
	pb "metrics/internal/proto/v1"
	pbModel "metrics/internal/proto/v1/model"
	"metrics/internal/service"

	"go.uber.org/zap"
)

type MetricServer struct {
	pb.UnimplementedMetricsServer
	metricService service.MetricService
	logger        *zap.SugaredLogger
}

func NewServer(svc service.MetricService, logger *zap.SugaredLogger) *MetricServer {
	return &MetricServer{
		metricService: svc,
		logger:        logger.With("component", "rpc MetricServer"),
	}
}

func (s *MetricServer) SendMetrics(ctx context.Context, req *pbModel.MetricsRequest) (*pbModel.MetricsResponse, error) {
	metrics := make([]service.MetricsUpdateRequest, 0, len(req.GetMetrics()))
	for _, m := range req.GetMetrics() {
		metrics = append(metrics, service.MetricsUpdateRequest{
			Delta: ptr(m.GetDelta()),
			Value: ptr(m.GetValue()),
			ID:    m.GetId(),
			MType: m.GetType(),
		})
	}

	err := s.metricService.UpdateMultiple(ctx, metrics)
	if err != nil {
		s.logger.Infow("service error", "error", err)
		return nil, fmt.Errorf("error update metrics: %w", err)
	}

	return &pbModel.MetricsResponse{
		Status: ptr("ok"),
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}
