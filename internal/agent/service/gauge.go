package service

import (
	"encoding/json"
	"fmt"
	"metrics/internal/agent/dto"

	"github.com/go-resty/resty/v2"
)

func SendMetric(client *resty.Client, request dto.MetricsGaugeUpdateRequest) error {
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
