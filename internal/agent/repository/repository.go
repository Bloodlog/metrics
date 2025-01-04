package repository

import (
	"math/rand"
	"runtime"
)

type Metric struct {
	Name  string
	Value uint64
}

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetMemoryMetrics() []Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := append([]Metric{},
		Metric{"Alloc", memStats.Alloc},
		Metric{"BuckHashSys", memStats.BuckHashSys},
		Metric{"Frees", memStats.Frees},
		Metric{"GCCPUFraction", uint64(memStats.GCCPUFraction)},
		Metric{"GCSys", memStats.GCSys},
		Metric{"HeapAlloc", memStats.HeapAlloc},
		Metric{"HeapIdle", memStats.HeapIdle},
		Metric{"HeapInuse", memStats.HeapInuse},
		Metric{"HeapObjects", memStats.HeapObjects},
		Metric{"HeapReleased", memStats.HeapReleased},
		Metric{"HeapSys", memStats.HeapSys},
		Metric{"LastGC", memStats.LastGC},
		Metric{"Lookups", memStats.Lookups},
		Metric{"MCacheInuse", memStats.MCacheInuse},
		Metric{"MCacheSys", memStats.MCacheSys},
		Metric{"MSpanInuse", memStats.MSpanInuse},
		Metric{"MSpanSys", memStats.MSpanSys},
		Metric{"Mallocs", memStats.Mallocs},
		Metric{"NextGC", memStats.NextGC},
		Metric{"NumForcedGC", uint64(memStats.NumForcedGC)},
		Metric{"NumGC", uint64(memStats.NumGC)},
		Metric{"OtherSys", memStats.OtherSys},
		Metric{"PauseTotalNs", memStats.PauseTotalNs},
		Metric{"StackInuse", memStats.StackInuse},
		Metric{"StackSys", memStats.StackSys},
		Metric{"Sys", memStats.Sys},
		Metric{"TotalAlloc", memStats.TotalAlloc},
		Metric{"RandomValue", rand.Uint64()},
	)

	return metrics
}
