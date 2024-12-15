package service

import (
	"errors"
	"fmt"
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

func Get(req MetricsGetRequest, storage *repository.MemStorage) (*MetricsResponse, error) {
	if req.MType == "counter" {
		counter, err := storage.GetCounter(req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get counter for ID %s: %w", req.ID, err)
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
			return nil, fmt.Errorf("failed to get gauge for ID %s: %w", req.ID, err)
		}
		gaugeValue := gauge
		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: &gaugeValue,
		}, nil
	}

	return nil, errors.New("type metric not found")
}

func Update(req MetricsUpdateRequest, storage *repository.MemStorage) (*MetricsResponse, error) {
	if req.MType == "counter" {
		deltaValue := uint64(*req.Delta)
		storage.SetCounter(req.ID, deltaValue)
		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: req.Delta,
			Value: nil,
		}, nil
	}

	if req.MType == "gauge" {
		value := *req.Value
		storage.SetGauge(req.ID, value)
		return &MetricsResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: nil,
			Value: req.Value,
		}, nil
	}

	return nil, errors.New("type metric not found")
}
