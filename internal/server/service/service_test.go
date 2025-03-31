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