package dto

// MetricsData представляет структуру для хранения данных метрик.
// Она содержит два поля: Gauges для хранения значений типа gauge и Counters для хранения значений типа counter.
type MetricsData struct {
	// Метрики.
	Gauges map[string]float64
	// Счетчики.
	Counters map[string]uint64
}
