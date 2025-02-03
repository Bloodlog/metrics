package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsGaugeUpdateRequest struct {
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func SendMetric(client *resty.Client, request MetricsGaugeUpdateRequest) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = client.R().
		SetBody(requestData).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %w", request.ID, err)
	}

	return nil
}
