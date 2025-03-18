package dto

// MetricsUpdateRequest Структура для обновления метрики.
type MetricsUpdateRequest struct {
	// Значение counter.
	Delta *int64 `json:"delta,omitempty"`
	// Значение gauge.
	Value *float64 `json:"value,omitempty"`
	// Имя метрики.
	ID string `json:"id"`
	// Тип метрики: counter или gauge.
	MType string `json:"type"`
}
