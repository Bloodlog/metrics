package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsUpdateRequest struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type MetricsUpdateRequests struct {
	Metrics []MetricsUpdateRequest `json:"metrics"`
}

func SendMetricsBatch(client *resty.Client, request MetricsUpdateRequests) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	compressedData, err := Compress(requestData)
	if err != nil {
		return fmt.Errorf("error compressing the data: %w", err)
	}

	_, err = client.R().
		SetBody(compressedData).
		Post("/updates")
	if err != nil {
		return fmt.Errorf("failed to send metric %w", err)
	}

	return nil
}
