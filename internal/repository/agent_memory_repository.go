package repository

import (
	"math/rand"
	"runtime"
)

type MemoryRepository struct {
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) GetMetrics() []Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := append([]Metric{},
		Metric{Name: "Alloc", Value: memStats.Alloc},
		Metric{Name: "BuckHashSys", Value: memStats.BuckHashSys},
		Metric{Name: "Frees", Value: memStats.Frees},
		Metric{Name: "GCCPUFraction", Value: uint64(memStats.GCCPUFraction)},
		Metric{Name: "GCSys", Value: memStats.GCSys},
		Metric{Name: "HeapAlloc", Value: memStats.HeapAlloc},
		Metric{Name: "HeapIdle", Value: memStats.HeapIdle},
		Metric{Name: "HeapInuse", Value: memStats.HeapInuse},
		Metric{Name: "HeapObjects", Value: memStats.HeapObjects},
		Metric{Name: "HeapReleased", Value: memStats.HeapReleased},
		Metric{Name: "HeapSys", Value: memStats.HeapSys},
		Metric{Name: "LastGC", Value: memStats.LastGC},
		Metric{Name: "Lookups", Value: memStats.Lookups},
		Metric{Name: "MCacheInuse", Value: memStats.MCacheInuse},
		Metric{Name: "MCacheSys", Value: memStats.MCacheSys},
		Metric{Name: "MSpanInuse", Value: memStats.MSpanInuse},
		Metric{Name: "MSpanSys", Value: memStats.MSpanSys},
		Metric{Name: "Mallocs", Value: memStats.Mallocs},
		Metric{Name: "NextGC", Value: memStats.NextGC},
		Metric{Name: "NumForcedGC", Value: uint64(memStats.NumForcedGC)},
		Metric{Name: "NumGC", Value: uint64(memStats.NumGC)},
		Metric{Name: "OtherSys", Value: memStats.OtherSys},
		Metric{Name: "PauseTotalNs", Value: memStats.PauseTotalNs},
		Metric{Name: "StackInuse", Value: memStats.StackInuse},
		Metric{Name: "StackSys", Value: memStats.StackSys},
		Metric{Name: "Sys", Value: memStats.Sys},
		Metric{Name: "TotalAlloc", Value: memStats.TotalAlloc},
		Metric{Name: "RandomValue", Value: rand.Uint64()},
	)

	return metrics
}
