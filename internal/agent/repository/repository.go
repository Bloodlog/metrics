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
	const metricsCount = 28

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := make([]Metric, metricsCount)

	metrics[0] = Metric{"Alloc", memStats.Alloc}
	metrics[1] = Metric{"BuckHashSys", memStats.BuckHashSys}
	metrics[2] = Metric{"Frees", memStats.Frees}
	metrics[3] = Metric{"GCCPUFraction", uint64(memStats.GCCPUFraction)}
	metrics[4] = Metric{"GCSys", memStats.GCSys}
	metrics[5] = Metric{"HeapAlloc", memStats.HeapAlloc}
	metrics[6] = Metric{"HeapIdle", memStats.HeapIdle}
	metrics[7] = Metric{"HeapInuse", memStats.HeapInuse}
	metrics[8] = Metric{"HeapObjects", memStats.HeapObjects}
	metrics[9] = Metric{"HeapReleased", memStats.HeapReleased}
	metrics[10] = Metric{"HeapSys", memStats.HeapSys}
	metrics[11] = Metric{"LastGC", memStats.LastGC}
	metrics[12] = Metric{"Lookups", memStats.Lookups}
	metrics[13] = Metric{"MCacheInuse", memStats.MCacheInuse}
	metrics[14] = Metric{"MCacheSys", memStats.MCacheSys}
	metrics[15] = Metric{"MSpanInuse", memStats.MSpanInuse}
	metrics[16] = Metric{"MSpanSys", memStats.MSpanSys}
	metrics[17] = Metric{"Mallocs", memStats.Mallocs}
	metrics[18] = Metric{"NextGC", memStats.NextGC}
	metrics[19] = Metric{"NumForcedGC", uint64(memStats.NumForcedGC)}
	metrics[20] = Metric{"NumGC", uint64(memStats.NumGC)}
	metrics[21] = Metric{"OtherSys", memStats.OtherSys}
	metrics[22] = Metric{"PauseTotalNs", memStats.PauseTotalNs}
	metrics[23] = Metric{"StackInuse", memStats.StackInuse}
	metrics[24] = Metric{"StackSys", memStats.StackSys}
	metrics[25] = Metric{"Sys", memStats.Sys}
	metrics[26] = Metric{"TotalAlloc", memStats.TotalAlloc}
	metrics[27] = Metric{"RandomValue", rand.Uint64()}

	return metrics
}
