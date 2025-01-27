package service

import (
	"context"
	"metrics/internal/server/repository"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	memStorage, _ := repository.NewMemStorage(ctx)
	logger := zap.NewNop()
	sugar := logger.Sugar()

	counterID := "testCounter"
	counterValue := uint64(42)
	_, err := memStorage.SetCounter(ctx, counterID, counterValue)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	gaugeID := "testGauge"
	gaugeValue := 123.45
	_, err = memStorage.SetGauge(ctx, gaugeID, gaugeValue)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	t.Run("Get counter metric", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    counterID,
			MType: "counter",
		}

		metricService := NewMetricService(sugar)
		resp, err := metricService.Get(ctx, req, memStorage)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, counterID, resp.ID)
		assert.Equal(t, "counter", resp.MType)
		assert.NotNil(t, resp.Delta)
		assert.Equal(t, int64(counterValue), *resp.Delta)
	})

	t.Run("Get gauge metric", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    gaugeID,
			MType: "gauge",
		}

		metricService := NewMetricService(sugar)
		resp, err := metricService.Get(ctx, req, memStorage)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, gaugeID, resp.ID)
		assert.Equal(t, "gauge", resp.MType)
		assert.NotNil(t, resp.Value)
		assert.Equal(t, gaugeValue, *resp.Value)
	})

	t.Run("Get non-existing metric", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    "unknownMetric",
			MType: "counter",
		}

		metricService := NewMetricService(sugar)
		resp, err := metricService.Get(ctx, req, memStorage)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid metric type", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    counterID,
			MType: "invalid",
		}

		metricService := NewMetricService(sugar)
		resp, err := metricService.Get(ctx, req, memStorage)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
