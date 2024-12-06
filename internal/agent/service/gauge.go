package service

import (
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
)

var (
	ErrSendMetric = errors.New("failed to send POST request PollCount")
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
		log.Printf("failed to send POST request PollCount: %v", err)
		return ErrSendMetric
	}

	return nil
}
