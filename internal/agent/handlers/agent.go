package handlers

import (
	"fmt"
	"metrics/internal/agent/clients"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
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
	Metrics   []repository.Metric
	PollCount int64
}
type Handlers struct {
	configs          *config.Config
	memoryRepository *repository.MemoryRepository
	systemRepository *repository.SystemRepository
	logger           *zap.SugaredLogger
	client           *resty.Client
	sendQueue        chan MetricsPayload
}

func NewHandlers(configs *config.Config, memoryRepository *repository.MemoryRepository, systemRepository *repository.SystemRepository, logger *zap.SugaredLogger) *Handlers {
	serverAddr := "http://" + net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)
	client := clients.NewClient(serverAddr, configs.Key, logger)

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
	for i := 0; i < h.configs.RateLimit; i++ {
		wg.Add(1)
		go h.worker(&wg)
	}

	var runtimeMetrics []repository.Metric
	var systemMetrics []repository.Metric
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
		allMetrics := append(runtimeMetrics, systemMetrics...)

		h.sendQueue <- MetricsPayload{
			Metrics:   allMetrics,
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
