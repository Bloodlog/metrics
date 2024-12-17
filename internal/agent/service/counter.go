package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricsCounterRequest struct {
	Delta *int64 `json:"delta,omitempty"`
	ID    string `json:"id"`
	MType string `json:"type"`
}

func SendIncrement(client *resty.Client, request MetricsCounterRequest) error {
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send increment: %w", err)
	}

	return nil
}
