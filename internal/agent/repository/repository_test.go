package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemoryMetrics(t *testing.T) {
	repo := NewRepository()

	metrics := repo.GetMemoryMetrics()

	assert.Equal(t, 28, len(metrics), "Количество метрик должно быть равно 28")

	for _, metric := range metrics {
		assert.NotEmpty(t, metric.Name, "Имя метрики не должно быть пустым")
	}

	expectedMetrics := map[string]bool{
		"Alloc":        true,
		"BuckHashSys":  true,
		"HeapAlloc":    true,
		"PauseTotalNs": true,
		"Sys":          true,
	}

	for _, metric := range metrics {
		delete(expectedMetrics, metric.Name)
	}

	assert.Empty(t, expectedMetrics, "Не все ожидаемые метрики найдены")
}
