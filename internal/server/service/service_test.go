package service

import (
	"metrics/internal/server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	memStorage := repository.NewMemStorage()

	counterID := "testCounter"
	counterValue := uint64(42)
	memStorage.SetCounter(counterID, counterValue)

	gaugeID := "testGauge"
	gaugeValue := 123.45
	memStorage.SetGauge(gaugeID, gaugeValue)

	t.Run("Get counter metric", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    counterID,
			MType: "counter",
		}

		resp, err := Get(req, memStorage)
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

		resp, err := Get(req, memStorage)
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

		resp, err := Get(req, memStorage)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid metric type", func(t *testing.T) {
		req := MetricsGetRequest{
			ID:    counterID,
			MType: "invalid",
		}

		resp, err := Get(req, memStorage)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
