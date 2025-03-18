package service

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/apperrors"
	"metrics/internal/server/dto"
	"metrics/internal/server/repository"

	"go.uber.org/zap"
)

type MetricService interface {
	Get(
		ctx context.Context,
		req dto.MetricsGetRequest,
	) (*dto.MetricsResponse, error)
	Update(
		ctx context.Context,
		req dto.MetricsUpdateRequest,
	) (*dto.MetricsResponse, error)
	UpdateMultiple(
		ctx context.Context,
		metrics []dto.MetricsUpdateRequest,
	) error
	GetMetrics(ctx context.Context) dto.MetricsData
}

type metricService struct {
	MetricRepository repository.MetricStorage
	logger           *zap.SugaredLogger
}

func NewMetricService(
	metricRepository repository.MetricStorage,
	logger *zap.SugaredLogger,
) MetricService {
	return &metricService{
		MetricRepository: metricRepository,
		logger:           logger,
	}
}

func (s *metricService) Get(
	ctx context.Context,
	req dto.MetricsGetRequest,
) (*dto.MetricsResponse, error) {
	if req.MType == "counter" {
		counter, err := s.MetricRepository.GetCounter(ctx, req.ID)
		if err != nil {
			return nil, apperrors.ErrMetricNotFound
		}

		counterValue := int64(counter)
		return &dto.MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: &counterValue,
			Value: nil,
		}, nil
	}

	if req.MType == "gauge" {
		gauge, err := s.MetricRepository.GetGauge(ctx, req.ID)
		if err != nil {
			return nil, apperrors.ErrMetricNotFound
		}
		gaugeValue := gauge
		return &dto.MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: &gaugeValue,
		}, nil
	}

	return nil, apperrors.ErrMetricNotFound
}

func (s *metricService) Update(
	ctx context.Context,
	req dto.MetricsUpdateRequest,
) (*dto.MetricsResponse, error) {
	if req.MType == "counter" {
		if req.Delta == nil {
			return nil, errors.New("delta field cannot be nil for counter type")
		}
		delta := *req.Delta
		deltaValue := uint64(delta)
		counter, err := s.MetricRepository.SetCounter(ctx, req.ID, deltaValue)
		if err != nil {
			return nil, errors.New("value cannot be save")
		}

		counterValue := int64(counter)
		return &dto.MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: &counterValue,
			Value: nil,
		}, nil
	}

	if req.MType == "gauge" {
		if req.Value == nil {
			return nil, errors.New("value field cannot be nil for gauge type")
		}
		value := *req.Value
		gauge, err := s.MetricRepository.SetGauge(ctx, req.ID, value)
		if err != nil {
			return nil, errors.New("value cannot be save")
		}

		gaugeValue := gauge

		return &dto.MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: &gaugeValue,
		}, nil
	}

	return nil, apperrors.ErrMetricNotFound
}

func (s *metricService) UpdateMultiple(
	ctx context.Context,
	metrics []dto.MetricsUpdateRequest,
) error {
	gauges := make(map[string]float64)
	counters := make(map[string]uint64)

	for _, metric := range metrics {
		if metric.Delta == nil && metric.Value == nil {
			continue
		}
		if metric.Delta != nil {
			counters[metric.ID] += uint64(*metric.Delta)
		}

		if metric.Value != nil {
			gauges[metric.ID] = *metric.Value
		}
	}

	err := s.MetricRepository.UpdateCounterAndGauges(ctx, counters, gauges)
	if err != nil {
		return fmt.Errorf("failed UpdateCounterAndGauges in service: %w", err)
	}

	return nil
}

func (s *metricService) GetMetrics(ctx context.Context) dto.MetricsData {
	gauges, _ := s.MetricRepository.Gauges(ctx)
	counters, _ := s.MetricRepository.Counters(ctx)

	data := dto.MetricsData{
		Gauges:   gauges,
		Counters: counters,
	}

	return data
}
