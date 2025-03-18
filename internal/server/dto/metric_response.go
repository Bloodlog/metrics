package dto

// MetricsResponse Структура для вывода метрики.
type MetricsResponse struct {
	// Значение counter.
	Delta *int64 `json:"delta,omitempty"`
	// Значение gauge.
	Value *float64 `json:"value,omitempty"`
	// Тип метрики: counter или gauge.
	ID string `json:"id"`
	// Имя метрики.
	MType string `json:"type"`
}
