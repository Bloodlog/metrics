package service

import (
	"context"
	"fmt"
	pb "metrics/internal/proto/v1"
	pbModel "metrics/internal/proto/v1/model"
)

type GRPCMetricSender struct {
	client pb.MetricsClient
}

func NewGRPCMetricSender(client pb.MetricsClient) *GRPCMetricSender {
	return &GRPCMetricSender{client: client}
}

func (s *GRPCMetricSender) SendIncrement(ctx context.Context, req AgentMetricsCounterRequest) error {
	metric := &pbModel.Metric{
		Id:    &req.ID,
		Type:  &req.MType,
		Delta: req.Delta,
	}

	return s.sendSingle(ctx, metric)
}

func (s *GRPCMetricSender) SendMetric(ctx context.Context, req AgentMetricsGaugeUpdateRequest) error {
	metric := &pbModel.Metric{
		Id:    &req.ID,
		Type:  &req.MType,
		Value: req.Value,
	}

	return s.sendSingle(ctx, metric)
}

func (s *GRPCMetricSender) sendSingle(ctx context.Context, metric *pbModel.Metric) error {
	request := &pbModel.MetricsRequest{
		Metrics: []*pbModel.Metric{metric},
	}

	_, err := s.client.SendMetrics(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to send metric via gRPC: %w", err)
	}

	return nil
}

func (s *GRPCMetricSender) SendMetricsBatch(ctx context.Context, req AgentMetricsUpdateRequests) error {
	grpcMetrics := make([]*pbModel.Metric, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		metric := &pbModel.Metric{
			Id:    &m.ID,
			Type:  &m.MType,
			Delta: m.Delta,
			Value: m.Value,
		}
		grpcMetrics = append(grpcMetrics, metric)
	}

	_, err := s.client.SendMetrics(ctx, &pbModel.MetricsRequest{
		Metrics: grpcMetrics,
	})
	if err != nil {
		return fmt.Errorf("failed to send batch via gRPC: %w", err)
	}

	return nil
}
