package dto

type MetricsCounterRequest struct {
	Delta *int64 `json:"delta,omitempty"`
	ID    string `json:"id"`
	MType string `json:"type"`
}
