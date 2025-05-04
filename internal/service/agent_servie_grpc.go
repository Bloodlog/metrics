package service

import (
	"context"
	"fmt"
	pb "metrics/internal/proto/v1"
	pbModel "metrics/internal/proto/v1/model"
	"time"
)

type GRPCMetricSender struct {
	client pb.MetricsClient
}

func NewGRPCMetricSender(client pb.MetricsClient) *GRPCMetricSender {
	return &GRPCMetricSender{client: client}
}

func (s *GRPCMetricSender) SendIncrement(req AgentMetricsCounterRequest) error {
	metric := &pbModel.Metric{
		Id:    &req.ID,
		Type:  &req.MType,
		Delta: req.Delta,
	}

	return s.sendSingle(metric)
}

func (s *GRPCMetricSender) SendMetric(req AgentMetricsGaugeUpdateRequest) error {
	metric := &pbModel.Metric{
		Id:    &req.ID,
		Type:  &req.MType,
		Value: req.Value,
	}

	return s.sendSingle(metric)
}

func (s *GRPCMetricSender) sendSingle(metric *pbModel.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	request := &pbModel.MetricsRequest{
		Metrics: []*pbModel.Metric{metric},
	}

	_, err := s.client.SendMetrics(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to send metric via gRPC: %w", err)
	}

	return nil
}

func (s *GRPCMetricSender) SendMetricsBatch(req AgentMetricsUpdateRequests) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.SendMetrics(ctx, &pbModel.MetricsRequest{
		Metrics: grpcMetrics,
	})
	if err != nil {
		return fmt.Errorf("failed to send batch via gRPC: %w", err)
	}

	return nil
}
