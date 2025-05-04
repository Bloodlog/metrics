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

func NewServer(service service.MetricService, logger *zap.SugaredLogger) *MetricServer {
	return &MetricServer{
		metricService: service,
		logger:        logger.With("component", "rpc MetricServer"),
	}
}

func (s *MetricServer) SendMetrics(ctx context.Context, req *pbModel.MetricsRequest) (*pbModel.MetricsResponse, error) {
	metrics := make([]service.MetricsUpdateRequest, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		metrics = append(metrics, service.MetricsUpdateRequest{
			Delta: m.Delta,
			Value: m.Value,
			ID:    getString(m.Id),
			MType: getString(m.Type),
		})
	}

	err := s.metricService.UpdateMultiple(ctx, metrics)
	if err != nil {
		s.logger.Infow("service error", "error", err)
		return nil, fmt.Errorf("error update metrics: %w", err)
	}

	return &pbModel.MetricsResponse{
		Status: strPtr("ok"),
	}, nil
}

func getString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func strPtr(s string) *string {
	return &s
}
