package repository

import (
	"math/rand"
	"metrics/internal/agent/dto"
	"runtime"
)

type MemoryRepository struct {
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) GetMetrics() []dto.Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := append([]dto.Metric{},
		dto.Metric{Name: "Alloc", Value: memStats.Alloc},
		dto.Metric{Name: "BuckHashSys", Value: memStats.BuckHashSys},
		dto.Metric{Name: "Frees", Value: memStats.Frees},
		dto.Metric{Name: "GCCPUFraction", Value: uint64(memStats.GCCPUFraction)},
		dto.Metric{Name: "GCSys", Value: memStats.GCSys},
		dto.Metric{Name: "HeapAlloc", Value: memStats.HeapAlloc},
		dto.Metric{Name: "HeapIdle", Value: memStats.HeapIdle},
		dto.Metric{Name: "HeapInuse", Value: memStats.HeapInuse},
		dto.Metric{Name: "HeapObjects", Value: memStats.HeapObjects},
		dto.Metric{Name: "HeapReleased", Value: memStats.HeapReleased},
		dto.Metric{Name: "HeapSys", Value: memStats.HeapSys},
		dto.Metric{Name: "LastGC", Value: memStats.LastGC},
		dto.Metric{Name: "Lookups", Value: memStats.Lookups},
		dto.Metric{Name: "MCacheInuse", Value: memStats.MCacheInuse},
		dto.Metric{Name: "MCacheSys", Value: memStats.MCacheSys},
		dto.Metric{Name: "MSpanInuse", Value: memStats.MSpanInuse},
		dto.Metric{Name: "MSpanSys", Value: memStats.MSpanSys},
		dto.Metric{Name: "Mallocs", Value: memStats.Mallocs},
		dto.Metric{Name: "NextGC", Value: memStats.NextGC},
		dto.Metric{Name: "NumForcedGC", Value: uint64(memStats.NumForcedGC)},
		dto.Metric{Name: "NumGC", Value: uint64(memStats.NumGC)},
		dto.Metric{Name: "OtherSys", Value: memStats.OtherSys},
		dto.Metric{Name: "PauseTotalNs", Value: memStats.PauseTotalNs},
		dto.Metric{Name: "StackInuse", Value: memStats.StackInuse},
		dto.Metric{Name: "StackSys", Value: memStats.StackSys},
		dto.Metric{Name: "Sys", Value: memStats.Sys},
		dto.Metric{Name: "TotalAlloc", Value: memStats.TotalAlloc},
		dto.Metric{Name: "RandomValue", Value: rand.Uint64()},
	)

	return metrics
}
