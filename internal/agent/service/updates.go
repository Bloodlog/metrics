package service

import (
	"encoding/json"
	"fmt"
	"metrics/internal/agent/dto"

	"github.com/go-resty/resty/v2"
)

func SendMetricsBatch(client *resty.Client, request dto.MetricsUpdateRequests) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = client.R().
		SetBody(requestData).
		Post("/updates")
	if err != nil {
		return fmt.Errorf("failed to send metric %w", err)
	}

	return nil
}
