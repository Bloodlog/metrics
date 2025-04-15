package handlers

import (
	"fmt"
	"metrics/internal/agent/clients"
	"metrics/internal/agent/dto"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"metrics/internal/config"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

const (
	nameCounter = "PollCount"
	typeCounter = "counter"
)

const typeMetricName = "gauge"

type MetricsPayload struct {
	Metrics   []dto.Metric
	PollCount int64
}
type Handlers struct {
	configs          *config.AgentConfig
	memoryRepository *repository.MemoryRepository
	systemRepository *repository.SystemRepository
	logger           *zap.SugaredLogger
	client           *resty.Client
	sendQueue        chan MetricsPayload
}

func NewHandlers(
	client *clients.Client,
	configs *config.AgentConfig,
	memoryRepository *repository.MemoryRepository,
	systemRepository *repository.SystemRepository,
	logger *zap.SugaredLogger,
) *Handlers {
	return &Handlers{
		configs:          configs,
		memoryRepository: memoryRepository,
		systemRepository: systemRepository,
		logger:           logger,
		client:           client.RestyClient,
	}
}

func (h *Handlers) Handle() error {
	pollTicker := time.NewTicker(time.Duration(h.configs.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(h.configs.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	h.sendQueue = make(chan MetricsPayload, h.configs.RateLimit)

	var wg sync.WaitGroup
	for range make([]struct{}, h.configs.RateLimit) {
		wg.Add(1)
		go h.worker(&wg)
	}

	var runtimeMetrics []dto.Metric
	var systemMetrics []dto.Metric
	counter := 0

	go func() {
		for range pollTicker.C {
			runtimeMetrics = h.memoryRepository.GetMetrics()
			counter++
		}
	}()

	go func() {
		for range pollTicker.C {
			systemMetrics = h.systemRepository.GetMetrics()
		}
	}()

	for range reportTicker.C {
		h.sendQueue <- MetricsPayload{
			Metrics:   append(runtimeMetrics, systemMetrics...),
			PollCount: int64(counter),
		}
	}

	close(h.sendQueue)
	wg.Wait()

	return nil
}

func (h *Handlers) worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for payload := range h.sendQueue {
		pollCount := payload.PollCount
		metrics := payload.Metrics

		var err error
		if h.configs.Batch {
			err = h.sendBatch(metrics, pollCount)
		} else {
			err = h.sendAPI(metrics, pollCount)
		}
		if err != nil {
			fmt.Printf("Failed to send metrics: %v\n", err)
		}
	}
}

func (h *Handlers) sendBatch(metrics []dto.Metric, counter int64) error {
	metricsRequests := dto.MetricsUpdateRequests{}
	metric := dto.MetricsUpdateRequest{
		Delta: &counter,
		ID:    nameCounter,
		MType: typeCounter,
	}
	metricsRequests.Metrics = append(metricsRequests.Metrics, metric)

	for _, metric := range metrics {
		valueFloat := float64(metric.Value)
		metric := dto.MetricsUpdateRequest{
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

func (h *Handlers) sendAPI(metrics []dto.Metric, counter int64) error {
	metricCounterRequest := dto.MetricsCounterRequest{
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

		MetricGaugeUpdateRequest := dto.MetricsGaugeUpdateRequest{
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
