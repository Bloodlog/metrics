package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetAndGetGauge(t *testing.T) {
	ctx := context.Background()
	const counterName string = "Allocate"
	const counterValue float64 = 1333.333
	ms, _ := NewMemStorage()

	_, err := ms.SetGauge(ctx, counterName, counterValue)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	value, err := ms.GetGauge(ctx, counterName)

	assert.NoError(t, err, "ошибка не должна быть")
	assert.Equal(t, counterValue, value, "значение gauge должно совпадать")
}

func TestMemStorage_GetGauge_NotFound(t *testing.T) {
	ctx := context.Background()
	ms, _ := NewMemStorage()

	const nameMetric = "nonexistent"
	_, err := ms.GetGauge(ctx, nameMetric)

	assert.Error(t, err, "ожидалась ошибка для несуществующего gauge")
	assert.EqualError(t, err, "gauge metric '"+nameMetric+"' not found")
}

func TestMemStorage_SetAndGetCounter(t *testing.T) {
	ctx := context.Background()
	ms, _ := NewMemStorage()

	_, err := ms.SetCounter(ctx, "requests", 5)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}
	_, err = ms.SetCounter(ctx, "requests", 10)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	value, err := ms.GetCounter(ctx, "requests")

	assert.NoError(t, err, "ошибка не должна быть")
	assert.Equal(t, uint64(15), value, "значение counter должно быть суммой")
}

func TestMemStorage_GetCounter_NotFound(t *testing.T) {
	ctx := context.Background()
	ms, _ := NewMemStorage()

	const nameMetric = "nonexistent"
	_, err := ms.GetCounter(ctx, nameMetric)

	assert.Error(t, err, "ожидалась ошибка для несуществующего counter")
	assert.EqualError(t, err, "counter metric '"+nameMetric+"' not found")
}

func TestMemStorage_UpdateCounterAndGauges(t *testing.T) {
	ctx := context.Background()
	ms, _ := NewMemStorage()

	counters := map[string]uint64{
		"counter1": 10,
		"counter2": 20,
	}
	gauges := map[string]float64{
		"gauge1": 1.1,
		"gauge2": 2.2,
	}

	err := ms.UpdateCounterAndGauges(ctx, counters, gauges)
	if err != nil {
		t.Errorf("Failed to UpdateCounterAndGauges: %v", err)
		return
	}

	for name, expected := range counters {
		actual, err := ms.GetCounter(ctx, name)
		assert.NoError(t, err, "ошибка не должна быть")
		assert.Equal(t, expected, actual)
	}

	for name, expected := range gauges {
		actual, err := ms.GetGauge(ctx, name)
		assert.NoError(t, err, "ошибка не должна быть")
		assert.Equal(t, expected, actual)
	}
}
