package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MetricSender interface {
	SendIncrement(req AgentMetricsCounterRequest) error
	SendMetric(req AgentMetricsGaugeUpdateRequest) error
	SendMetricsBatch(req AgentMetricsUpdateRequests) error
}

type HTTPMetricSender struct {
	client *resty.Client
}

func NewHTTPMetricSender(client *resty.Client) *HTTPMetricSender {
	return &HTTPMetricSender{client: client}
}

type AgentMetricsCounterRequest struct {
	Delta *int64 `json:"delta,omitempty"`
	ID    string `json:"id"`
	MType string `json:"type"`
}

type AgentMetricsGaugeUpdateRequest struct {
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type AgentMetricsUpdateRequest struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type AgentMetricsUpdateRequests struct {
	Metrics []AgentMetricsUpdateRequest `json:"metrics"`
}

func (s *HTTPMetricSender) SendIncrement(request AgentMetricsCounterRequest) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = s.client.R().
		SetBody(requestData).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send increment: %w", err)
	}

	return nil
}

func (s *HTTPMetricSender) SendMetric(request AgentMetricsGaugeUpdateRequest) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = s.client.R().
		SetBody(requestData).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %w", request.ID, err)
	}

	return nil
}

func (s *HTTPMetricSender) SendMetricsBatch(request AgentMetricsUpdateRequests) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error serializing the structure: %w", err)
	}

	_, err = s.client.R().
		SetBody(requestData).
		Post("/updates")
	if err != nil {
		return fmt.Errorf("failed to send metric %w", err)
	}

	return nil
}
