package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsCounterRequest struct {
	Delta *int64 `json:"delta,omitempty"`
	ID    string `json:"id"`
	MType string `json:"type"`
}

func SendIncrement(client *resty.Client, request MetricsCounterRequest) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = client.R().
		SetBody(requestData).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send increment: %w", err)
	}

	return nil
}
