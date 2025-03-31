package service

import (
	"context"
	"metrics/internal/server/dto"
	"metrics/internal/server/repository"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	memStorage, _ := repository.NewMemStorage()
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
		req := dto.MetricsGetRequest{
			ID:    counterID,
			MType: "counter",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Get(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, counterID, resp.ID)
		assert.Equal(t, "counter", resp.MType)
		assert.NotNil(t, resp.Delta)
		assert.Equal(t, int64(counterValue), *resp.Delta)
	})

	t.Run("Get gauge metric", func(t *testing.T) {
		req := dto.MetricsGetRequest{
			ID:    gaugeID,
			MType: "gauge",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Get(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, gaugeID, resp.ID)
		assert.Equal(t, "gauge", resp.MType)
		assert.NotNil(t, resp.Value)
		assert.Equal(t, gaugeValue, *resp.Value)
	})

	t.Run("Get non-existing metric", func(t *testing.T) {
		req := dto.MetricsGetRequest{
			ID:    "unknownMetric",
			MType: "counter",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid metric type", func(t *testing.T) {
		req := dto.MetricsGetRequest{
			ID:    counterID,
			MType: "invalid",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	memStorage, _ := repository.NewMemStorage()
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
		t.Errorf("Failed to SetGauge: %v", err)
		return
	}

	t.Run("Update counter metric", func(t *testing.T) {
		req := dto.MetricsUpdateRequest{
			Delta: new(int64),
			ID:    counterID,
			MType: "counter",
		}
		*req.Delta = 10

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Update(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, counterID, resp.ID)
		assert.Equal(t, "counter", resp.MType)
		assert.NotNil(t, resp.Delta)
		assert.Equal(t, int64(52), *resp.Delta)
	})

	t.Run("Update gauge metric", func(t *testing.T) {
		req := dto.MetricsUpdateRequest{
			Value: new(float64),
			ID:    gaugeID,
			MType: "gauge",
		}
		*req.Value = 150.5

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Update(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, gaugeID, resp.ID)
		assert.Equal(t, "gauge", resp.MType)
		assert.NotNil(t, resp.Value)
		assert.Equal(t, 150.5, *resp.Value)
	})

	t.Run("Update counter metric with nil Delta", func(t *testing.T) {
		req := dto.MetricsUpdateRequest{
			Delta: nil,
			ID:    counterID,
			MType: "counter",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Update(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "delta field cannot be nil for counter type", err.Error())
	})

	t.Run("Update gauge metric with nil Value", func(t *testing.T) {
		req := dto.MetricsUpdateRequest{
			Value: nil,
			ID:    gaugeID,
			MType: "gauge",
		}

		metricService := NewMetricService(memStorage, sugar)
		resp, err := metricService.Update(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "value field cannot be nil for gauge type", err.Error())
	})
}

func TestGetMetrics(t *testing.T) {
	ctx := context.Background()
	memStorage, _ := repository.NewMemStorage()
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
		t.Errorf("Failed to SetGauge: %v", err)
		return
	}

	t.Run("Get metrics", func(t *testing.T) {
		metricService := NewMetricService(memStorage, sugar)
		resp := metricService.GetMetrics(ctx)

		assert.Equal(t, resp.Counters[counterID], counterValue)
		assert.Equal(t, resp.Gauges[gaugeID], gaugeValue)
	})
}

func TestUpdateMultiple(t *testing.T) {
	ctx := context.Background()
	memStorage, _ := repository.NewMemStorage()
	logger := zap.NewNop()
	sugar := logger.Sugar()
	metricService := NewMetricService(memStorage, sugar)

	counterID := "testCounter"
	gaugeID := "testGauge"

	_, _ = memStorage.SetCounter(ctx, counterID, 42)
	_, _ = memStorage.SetGauge(ctx, gaugeID, 123.45)

	t.Run("Update multiple metrics", func(t *testing.T) {
		metrics := []dto.MetricsUpdateRequest{
			{
				ID:    counterID,
				MType: "counter",
				Delta: new(int64),
			},
			{
				ID:    gaugeID,
				MType: "gauge",
				Value: new(float64),
			},
		}
		*metrics[0].Delta = 10
		*metrics[1].Value = 150.5

		err := metricService.UpdateMultiple(ctx, metrics)
		assert.NoError(t, err)

		updatedCounter, _ := memStorage.GetCounter(ctx, counterID)
		updatedGauge, _ := memStorage.GetGauge(ctx, gaugeID)

		assert.Equal(t, uint64(52), updatedCounter)
		assert.Equal(t, 150.5, updatedGauge)
	})

	t.Run("Update with nil values", func(t *testing.T) {
		metrics := []dto.MetricsUpdateRequest{
			{
				ID:    counterID,
				MType: "counter",
				Delta: nil,
			},
			{
				ID:    gaugeID,
				MType: "gauge",
				Value: nil,
			},
		}

		err := metricService.UpdateMultiple(ctx, metrics)
		assert.NoError(t, err)
	})
}
