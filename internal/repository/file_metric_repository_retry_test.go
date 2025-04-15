package repository

import (
	"context"
	"metrics/internal/config"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func setupFileRetryStorageWrapper(t *testing.T) *FileRetryStorageWrapper {
	t.Helper()
	const tempFileNamePattern = "metrics_test_retry.json"
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

	_, err := fs.GetGauge(ctx, "test_gauge")
	assert.NoError(t, err)
}

func TestGetCounterRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
	ctx := context.Background()
	_, _ = fs.SetCounter(ctx, "test_counter", 5)

	_, err := fs.GetCounter(ctx, "test_counter")
	assert.NoError(t, err)
}

func TestGaugesRetry(t *testing.T) {
	fs := setupFileRetryStorageWrapper(t)
	ctx := context.Background()
	_, _ = fs.SetGauge(ctx, "gauge1", 10.1)
	_, _ = fs.SetGauge(ctx, "gauge2", 20.2)

	_, err := fs.Gauges(ctx)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
