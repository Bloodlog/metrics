package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func SendMetric(client *resty.Client, name string, value string) error {
	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"metricName":  name,
			"metricValue": value,
		}).
		Post("/update/gauge/{metricName}/{metricValue}")
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %w", name, err)
	}

	return nil
}
