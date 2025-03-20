package dto

// MetricsGetRequest Структура для запроса получения метрики.
type MetricsGetRequest struct {
	// Имя метрики.
	ID string `json:"id"`
	// Тип метрики: counter или gauge.
	MType string `json:"type"`
}
