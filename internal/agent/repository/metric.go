package repository

type Metric struct {
	Name  string
	Value uint64
}

type MetricsRepository interface {
	GetMetrics() []Metric
}
