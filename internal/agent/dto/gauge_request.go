package dto

type MetricsGaugeUpdateRequest struct {
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}
