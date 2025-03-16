package repository

// Metric хранит информацию о метриках.
type Metric struct {
	Name  string
	Value uint64
}

type MetricsRepository interface {
	GetMetrics() []Metric
}
