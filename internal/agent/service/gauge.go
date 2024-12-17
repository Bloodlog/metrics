package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsUpdateRequest struct {
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func SendMetric(client *resty.Client, request MetricsUpdateRequest) error {
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %w", request.ID, err)
	}

	return nil
}
