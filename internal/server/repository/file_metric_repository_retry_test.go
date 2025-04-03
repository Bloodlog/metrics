package repository

import (
	"context"
	"metrics/internal/server/config"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func setupFileRetryStorageWrapper(t *testing.T) *FileRetryStorageWrapper {
	const tempFileNamePattern = "metrics_test_retry.json"
	tempFile, err := os.CreateTemp("", tempFileNamePattern)
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return nil
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	cfg := &config.Config{
		FileStoragePath: tempFileNamePattern,
		StoreInterval:   1,
		Restore:         false,
	}
	logger := zap.NewNop().Sugar()
	fs, err := NewRetryFileStorage(context.Background(), cfg, logger)
	if err != nil {
		t.Fatalf("failed to create FileStorageWrapper: %v", err)
	}
	return fs
}

func TestGetGaugeRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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

func TestGetCounterRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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

func TestGaugesRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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

func TestCountersRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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

func TestUpdateCounterAndGaugesRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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

func TestAutoSaveRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
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
