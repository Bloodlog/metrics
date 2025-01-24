package repository

import (
	"context"
	"metrics/internal/server/config"
	"os"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestSaveToFileAndLoadFromFile(t *testing.T) {
	ctx := context.Background()
	const tempFileNamePattern = "metrics_test_*.json"

	tempFile, err := os.CreateTemp("", tempFileNamePattern)
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	cfg := &config.Config{
		FileStoragePath: tempFile.Name(),
		Restore:         false,
		StoreInterval:   0,
	}

	logger := zap.NewNop()
	sugar := logger.Sugar()

	expectedGauges := map[string]float64{
		"gauge1": 123.45,
		"gauge2": 678.90,
	}
	expectedCounters := map[string]uint64{
		"counter1": 100,
		"counter2": 200,
	}

	fileWrapper, err := NewFileStorageWrapper(ctx, cfg, sugar)
	if err != nil {
		t.Errorf("Failed to create FileStorageWrapper: %v", err)
		return
	}

	for name, value := range expectedGauges {
		if _, err := fileWrapper.SetGauge(ctx, name, value); err != nil {
			t.Errorf("Failed to set gauge '%s': %v", name, err)
			return
		}
	}
	for name, value := range expectedCounters {
		if _, err := fileWrapper.SetCounter(ctx, name, value); err != nil {
			t.Errorf("Failed to set counter '%s': %v", name, err)
			return
		}
	}

	if err := fileWrapper.saveToFile(ctx); err != nil {
		t.Errorf("saveToFile failed: %v", err)
		return
	}

	newStorage := &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}
	fileWrapper.storage = newStorage

	if err := fileWrapper.loadFromFile(ctx); err != nil {
		t.Errorf("loadFromFile failed: %v", err)
		return
	}

	for name, expectedValue := range expectedGauges {
		actualValue, err := newStorage.GetGauge(ctx, name)
		if err != nil {
			t.Errorf("Failed to get gauge '%s': %v", name, err)
			return
		}
		if actualValue != expectedValue {
			t.Errorf("Gauge '%s': expected %v, got %v", name, expectedValue, actualValue)
			return
		}
	}

	for name, expectedValue := range expectedCounters {
		actualValue, err := newStorage.GetCounter(ctx, name)
		if err != nil {
			t.Errorf("Failed to get counter '%s': %v", name, err)
			return
		}
		if actualValue != expectedValue {
			t.Errorf("Counter '%s': expected %v, got %v", name, expectedValue, actualValue)
			return
		}
	}
}
