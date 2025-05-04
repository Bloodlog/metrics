package service

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/repository"

	"go.uber.org/zap"
)

var ErrMetricNotFound = errors.New("metric not found")

// MetricsUpdateRequests Структура, содержащая данные метрик.
type MetricsUpdateRequests struct {
	Metrics []MetricsUpdateRequest `json:"metrics"`
}

// MetricsGetRequest Структура для запроса получения метрики.
type MetricsGetRequest struct {
	// Имя метрики.
	ID string `json:"id"`
	// Тип метрики: counter или gauge.
	MType string `json:"type"`
}

// MetricsResponse Структура для вывода метрики.
type MetricsResponse struct {
	// Значение counter.
	Delta *int64 `json:"delta,omitempty"`
	// Значение gauge.
	Value *float64 `json:"value,omitempty"`
	// Тип метрики: counter или gauge.
	ID string `json:"id"`
	// Имя метрики.
	MType string `json:"type"`
}

// MetricsData представляет структуру для хранения данных метрик.
// Она содержит два поля: Gauges для хранения значений типа gauge и Counters для хранения значений типа counter.
type MetricsData struct {
	// Метрики.
	Gauges map[string]float64
	// Счетчики.
	Counters map[string]uint64
}

// MetricsUpdateRequest Структура для обновления метрики.
type MetricsUpdateRequest struct {
	// Значение counter.
	Delta *int64 `json:"delta,omitempty"`
	// Значение gauge.
	Value *float64 `json:"value,omitempty"`
	// Имя метрики.
	ID string `json:"id"`
	// Тип метрики: counter или gauge.
	MType string `json:"type"`
}

type MetricService interface {
	Get(
		ctx context.Context,
		req MetricsGetRequest,
	) (*MetricsResponse, error)
	Update(
		ctx context.Context,
		req MetricsUpdateRequest,
	) (*MetricsResponse, error)
	UpdateMultiple(
		ctx context.Context,
		metrics []MetricsUpdateRequest,
	) error
	GetMetrics(ctx context.Context) MetricsData
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
	req MetricsGetRequest,
) (*MetricsResponse, error) {
	if req.MType == "counter" {
		counter, err := s.MetricRepository.GetCounter(ctx, req.ID)
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
		gauge, err := s.MetricRepository.GetGauge(ctx, req.ID)
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

func (s *metricService) Update(
	ctx context.Context,
	req MetricsUpdateRequest,
) (*MetricsResponse, error) {
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
		gauge, err := s.MetricRepository.SetGauge(ctx, req.ID, value)
		if err != nil {
			return nil, errors.New("value cannot be save")
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

func (s *metricService) UpdateMultiple(
	ctx context.Context,
	metrics []MetricsUpdateRequest,
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

func (s *metricService) GetMetrics(ctx context.Context) MetricsData {
	gauges, _ := s.MetricRepository.Gauges(ctx)
	counters, _ := s.MetricRepository.Counters(ctx)

	data := MetricsData{
		Gauges:   gauges,
		Counters: counters,
	}

	return data
}
