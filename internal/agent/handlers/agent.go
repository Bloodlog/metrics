package handlers

import (
	"fmt"
	"metrics/internal/agent/clients"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

const (
	nameCounter = "PollCount"
	typeCounter = "counter"
)

const typeMetricName = "gauge"

type Handlers struct {
	configs *config.Config
	storage *repository.Repository
	logger  *zap.SugaredLogger
	client  *resty.Client
}

func NewHandlers(configs *config.Config, storage *repository.Repository, logger *zap.SugaredLogger) *Handlers {
	serverAddr := "http://" + net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)
	client := clients.CreateClient(serverAddr, logger)

	return &Handlers{
		configs: configs,
		storage: storage,
		logger:  logger,
		client:  client,
	}
}

func (h *Handlers) Handle() error {
	pollTicker := time.NewTicker(time.Duration(h.configs.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(h.configs.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	var metrics []repository.Metric
	counter := 0

	for {
		select {
		case <-pollTicker.C:
			metrics = h.storage.GetMemoryMetrics()
			counter++

		case <-reportTicker.C:
			delta := int64(counter)
			var err error
			if h.configs.Batch {
				err = h.sendBatch(metrics, delta)
			} else {
				err = h.sendAPI(metrics, delta)
			}
			if err != nil {
				return fmt.Errorf("failed to send metric to server: %w", err)
			}
			counter = 0
		}
	}
}

func (h *Handlers) sendBatch(metrics []repository.Metric, counter int64) error {
	metricsRequests := service.MetricsUpdateRequests{}
	metric := service.MetricsUpdateRequest{
		Delta: &counter,
		ID:    nameCounter,
		MType: typeCounter,
	}
	metricsRequests.Metrics = append(metricsRequests.Metrics, metric)

	for _, metric := range metrics {
		valueFloat := float64(metric.Value)
		metric := service.MetricsUpdateRequest{
			Value: &valueFloat,
			ID:    metric.Name,
			MType: typeMetricName,
		}
		metricsRequests.Metrics = append(metricsRequests.Metrics, metric)
	}

	err := service.SendMetricsBatch(h.client, metricsRequests)
	if err != nil {
		return fmt.Errorf("failed to send metric to server: %w", err)
	}

	return nil
}

func (h *Handlers) sendAPI(metrics []repository.Metric, counter int64) error {
	metricCounterRequest := service.MetricsCounterRequest{
		Delta: &counter,
		ID:    nameCounter,
		MType: typeCounter,
	}

	err := service.SendIncrement(h.client, metricCounterRequest)
	if err != nil {
		return fmt.Errorf("failed to send Increment %d to server: %w", counter, err)
	}
	counter = 0

	for _, metric := range metrics {
		valueFloat := float64(metric.Value)

		MetricGaugeUpdateRequest := service.MetricsGaugeUpdateRequest{
			Value: &valueFloat,
			ID:    metric.Name,
			MType: typeMetricName,
		}

		err = service.SendMetric(h.client, MetricGaugeUpdateRequest)
		if err != nil {
			return fmt.Errorf("failed to send metric %s to server: %w", metric.Name, err)
		}
	}

	return nil
}
