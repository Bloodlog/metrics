package dto

type MetricsUpdateRequest struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type MetricsUpdateRequests struct {
	Metrics []MetricsUpdateRequest `json:"metrics"`
}
