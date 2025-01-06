package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsUpdateRequest struct {
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func SendMetric(client *resty.Client, request MetricsUpdateRequest) error {
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
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %w", request.ID, err)
	}

	return nil
}
