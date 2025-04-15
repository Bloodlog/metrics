package repository

import (
	"context"
	"metrics/internal/config"
	"os"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func createTempFile(t *testing.T, pattern string) *os.File {
	t.Helper()
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tempFile
}

func setupTestFileStorage(t *testing.T) *FileStorageWrapper {
	t.Helper()
	const tempFileNamePattern = "metrics_test.json"
	tempFile := createTempFile(t, tempFileNamePattern)
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	cfg := &config.ServerConfig{
		FileStoragePath: tempFileNamePattern,
		StoreInterval:   1,
		Restore:         false,
	}
	logger := zap.NewNop().Sugar()
	fs, err := NewFileStorageWrapper(context.Background(), cfg, logger)
	if err != nil {
		t.Fatalf("failed to create FileStorageWrapper: %v", err)
	}
	return fs
}

func TestSaveToFileAndLoadFromFile(t *testing.T) {
	ctx := context.Background()
	const tempFileNamePattern = "metrics_test_*.json"
	tempFile := createTempFile(t, tempFileNamePattern)
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	cfg := &config.ServerConfig{
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

	if err = fileWrapper.saveToFile(ctx); err != nil {
		t.Errorf("saveToFile failed: %v", err)
		return
	}

	newStorage := &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}
	fileWrapper.storage = newStorage

	if err = fileWrapper.loadFromFile(ctx); err != nil {
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

func TestGetGauge(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()
	_, _ = fs.SetGauge(ctx, "test_gauge", 42.42)

	value, err := fs.GetGauge(ctx, "test_gauge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 42.42 {
		t.Errorf("expected 42.42, got %v", value)
	}
}

func TestGetCounter(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()
	_, _ = fs.SetCounter(ctx, "test_counter", 5)

	value, err := fs.GetCounter(ctx, "test_counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 5 {
		t.Errorf("expected 5, got %v", value)
	}
}

func TestGauges(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()
	_, _ = fs.SetGauge(ctx, "gauge1", 10.1)
	_, _ = fs.SetGauge(ctx, "gauge2", 20.2)

	gauges, err := fs.Gauges(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gauges) != 2 {
		t.Errorf("expected 2 gauges, got %d", len(gauges))
	}
}

func TestCounters(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()
	_, _ = fs.SetCounter(ctx, "counter1", 10)
	_, _ = fs.SetCounter(ctx, "counter2", 20)

	counters, err := fs.Counters(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(counters) != 2 {
		t.Errorf("expected 2 counters, got %d", len(counters))
	}
}

func TestUpdateCounterAndGauges(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()
	counters := map[string]uint64{"counter1": 15}
	gauges := map[string]float64{"gauge1": 25.5}
	err := fs.UpdateCounterAndGauges(ctx, counters, gauges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	valCounter, _ := fs.GetCounter(ctx, "counter1")
	if valCounter != 15 {
		t.Errorf("expected 15, got %d", valCounter)
	}
	valGauge, _ := fs.GetGauge(ctx, "gauge1")
	if valGauge != 25.5 {
		t.Errorf("expected 25.5, got %v", valGauge)
	}
}

func TestAutoSave(t *testing.T) {
	fs := setupTestFileStorage(t)
	ctx := context.Background()

	_, _ = fs.SetGauge(ctx, "test_gauge", 100.5)
	_, _ = fs.SetCounter(ctx, "test_counter", 500)

	time.Sleep(2 * time.Second)

	_, err := os.Stat(fs.cfg.FileStoragePath)
	if os.IsNotExist(err) {
		t.Fatal("autosave did not create the file")
	}
	_ = os.Remove(fs.cfg.FileStoragePath)
}
