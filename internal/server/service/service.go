package service

import (
	"errors"
	"metrics/internal/server/repository"
)

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

func Get(req MetricsGetRequest, storage *repository.MemStorage) (*MetricsResponse, error) {
	if req.MType == "counter" {
		counter, err := storage.GetCounter(req.ID)
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
		gauge, err := storage.GetGauge(req.ID)
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

func Update(req MetricsUpdateRequest, storage *repository.MemStorage) (*MetricsResponse, error) {
	if req.MType == "counter" {
		if req.Delta == nil {
			return nil, errors.New("delta field cannot be nil for counter type")
		}
		delta := *req.Delta
		deltaValue := uint64(delta)
		storage.SetCounter(req.ID, deltaValue)

		counter, err := storage.GetCounter(req.ID)
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
		storage.SetGauge(req.ID, value)

		gauge, err := storage.GetGauge(req.ID)
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
