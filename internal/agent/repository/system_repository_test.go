package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemMetrics(t *testing.T) {
	repo := NewSystemRepository()

	metrics := repo.GetMetrics()

	assert.Equal(t, 12, len(metrics), "Количество метрик должно быть равно 12")

	for _, metric := range metrics {
		assert.NotEmpty(t, metric.Name, "Имя метрики не должно быть пустым")
	}

	expectedMetrics := map[string]bool{
		"TotalMemory": true,
		"FreeMemory":  true,
	}

	for _, metric := range metrics {
		delete(expectedMetrics, metric.Name)
	}

	assert.Empty(t, expectedMetrics, "Не все ожидаемые метрики найдены")
}
