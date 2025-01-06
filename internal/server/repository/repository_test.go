package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetAndGetGauge(t *testing.T) {
	const counterName string = "Allocate"
	const counterValue float64 = 1333.333
	ms := NewMemStorage()

	err := ms.SetGauge(counterName, counterValue)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	value, err := ms.GetGauge(counterName)

	assert.NoError(t, err, "ошибка не должна быть")
	assert.Equal(t, counterValue, value, "значение gauge должно совпадать")
}

func TestMemStorage_GetGauge_NotFound(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetGauge("nonexistent")

	assert.Error(t, err, "ожидалась ошибка для несуществующего gauge")
	assert.EqualError(t, err, "gauge metric not found")
}

func TestMemStorage_SetAndGetCounter(t *testing.T) {
	ms := NewMemStorage()

	err := ms.SetCounter("requests", 5)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}
	err = ms.SetCounter("requests", 10)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	value, err := ms.GetCounter("requests")

	assert.NoError(t, err, "ошибка не должна быть")
	assert.Equal(t, uint64(15), value, "значение counter должно быть суммой")
}

func TestMemStorage_GetCounter_NotFound(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetCounter("nonexistent")

	assert.Error(t, err, "ожидалась ошибка для несуществующего counter")
	assert.EqualError(t, err, "counter metric not found")
}
