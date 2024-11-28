package repository

import "errors"

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MetricStorage interface {
	SetGauge(name string, value float64)
	GetGauge(name string) (float64, error)
	SetCounter(name string, value uint64)
	GetCounter(name string) (uint64, error)
}

type MemStorage struct {
	gauges   map[string]float64 // Хранилище для метрик типа gauge
	counters map[string]uint64  // Хранилище для метрик типа counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
	}
}

func (ms *MemStorage) SetGauge(name string, value float64) {
	ms.gauges[name] = value
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	value, exists := ms.gauges[name]
	if !exists {
		return 0, errors.New("gauge metric not found")
	}
	return value, nil
}

func (ms *MemStorage) SetCounter(name string, value uint64) {
	ms.counters[name] += value
}

func (ms *MemStorage) GetCounter(name string) (uint64, error) {
	value, exists := ms.counters[name]
	if !exists {
		return 0, errors.New("counter metric not found")
	}
	return value, nil
}
