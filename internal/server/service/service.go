package service

import (
	"context"
	"errors"
	"metrics/internal/server/repository"

	"go.uber.org/zap"
)

type MetricService struct {
	logger *zap.SugaredLogger
}

func NewMetricService(logger *zap.SugaredLogger) *MetricService {
	return &MetricService{logger: logger}
}

type MetricsGetRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type MetricsUpdateRequest struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type MetricsResponse struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

var ErrMetricNotFound = errors.New("metric not found")

func (s *MetricService) Get(
	ctx context.Context,
	req MetricsGetRequest,
	storage repository.MetricStorage) (*MetricsResponse, error) {
	if req.MType == "counter" {
		counter, err := storage.GetCounter(ctx, req.ID)
		if err != nil {
			return nil, ErrMetricNotFound
		}

		counterValue := int64(counter)
		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: &counterValue,
			Value: nil,
		}, nil
	}

	if req.MType == "gauge" {
		gauge, err := storage.GetGauge(ctx, req.ID)
		if err != nil {
			return nil, ErrMetricNotFound
		}
		gaugeValue := gauge
		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: &gaugeValue,
		}, nil
	}

	return nil, ErrMetricNotFound
}

func (s *MetricService) Update(
	ctx context.Context,
	req MetricsUpdateRequest,
	storage repository.MetricStorage) (*MetricsResponse, error) {
	handlerLogger := s.logger.With("service", "Update")
	if req.MType == "counter" {
		if req.Delta == nil {
			return nil, errors.New("delta field cannot be nil for counter type")
		}
		delta := *req.Delta
		deltaValue := uint64(delta)
		err := storage.SetCounter(ctx, req.ID, deltaValue)
		if err != nil {
			handlerLogger.Infow("error: value cannot be save", "error", err)
			return nil, errors.New("value cannot be save")
		}

		counter, err := storage.GetCounter(ctx, req.ID)
		if err != nil {
			return nil, ErrMetricNotFound
		}

		counterValue := int64(counter)
		return &MetricsResponse{
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
		err := storage.SetGauge(ctx, req.ID, value)
		if err != nil {
			handlerLogger.Infow("error: value cannot be save", "error", err)
			return nil, errors.New("value cannot be save")
		}

		gauge, err := storage.GetGauge(ctx, req.ID)
		if err != nil {
			return nil, ErrMetricNotFound
		}
		gaugeValue := gauge

		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: &gaugeValue,
		}, nil
	}

	return nil, ErrMetricNotFound
}
