package handlers

import (
	"context"
	"fmt"
	"metrics/internal/config"
	"metrics/internal/repository"
	"metrics/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
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

type AgentHandler struct {
	configs          *config.AgentConfig
	memoryRepository *repository.MemoryRepository
	systemRepository *repository.SystemRepository
	agentService     service.MetricSender
	logger           *zap.SugaredLogger
	sendQueue        chan MetricsPayload
}

func NewAgentHandler(
	configs *config.AgentConfig,
	memoryRepository *repository.MemoryRepository,
	systemRepository *repository.SystemRepository,
	metricService service.MetricSender,
	logger *zap.SugaredLogger,
) *AgentHandler {
	return &AgentHandler{
		configs:          configs,
		memoryRepository: memoryRepository,
		systemRepository: systemRepository,
		agentService:     metricService,
		logger:           logger,
	}
}

func (h *AgentHandler) Handle() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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

	var runtimeMetrics []repository.Metric
	var systemMetrics []repository.Metric
	var counter int64 = 0

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				runtimeMetrics = h.memoryRepository.GetMetrics()
				counter++
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				systemMetrics = h.systemRepository.GetMetrics()
			}
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-reportTicker.C:
			h.sendQueue <- MetricsPayload{
				Metrics:   append(runtimeMetrics, systemMetrics...),
				PollCount: counter,
			}
		case sig := <-sigCh:
			h.logger.Infof("Received signal: %s, shutting down...", sig)
			cancel()
			break loop
		}
	}
	h.sendQueue <- MetricsPayload{
		Metrics:   append(runtimeMetrics, systemMetrics...),
		PollCount: counter,
	}

	close(h.sendQueue)
	wg.Wait()

	h.logger.Info("Agent gracefully shut down")
	return nil
}

func (h *AgentHandler) worker(wg *sync.WaitGroup) {
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

func (h *AgentHandler) sendBatch(metrics []repository.Metric, counter int64) error {
	metricsRequests := service.AgentMetricsUpdateRequests{}
	metric := service.AgentMetricsUpdateRequest{
		Delta: &counter,
		ID:    nameCounter,
		MType: typeCounter,
	}
	metricsRequests.Metrics = append(metricsRequests.Metrics, metric)

	for _, metric := range metrics {
		valueFloat := float64(metric.Value)
		metric := service.AgentMetricsUpdateRequest{
			Value: &valueFloat,
			ID:    metric.Name,
			MType: typeMetricName,
		}
		metricsRequests.Metrics = append(metricsRequests.Metrics, metric)
	}

	err := h.agentService.SendMetricsBatch(metricsRequests)
	if err != nil {
		return fmt.Errorf("failed to send metric to server: %w", err)
	}

	return nil
}

func (h *AgentHandler) sendAPI(metrics []repository.Metric, counter int64) error {
	metricCounterRequest := service.AgentMetricsCounterRequest{
		Delta: &counter,
		ID:    nameCounter,
		MType: typeCounter,
	}

	err := h.agentService.SendIncrement(metricCounterRequest)
	if err != nil {
		return fmt.Errorf("failed to send Increment %d to server: %w", counter, err)
	}
	counter = 0

	for _, metric := range metrics {
		valueFloat := float64(metric.Value)

		MetricGaugeUpdateRequest := service.AgentMetricsGaugeUpdateRequest{
			Value: &valueFloat,
			ID:    metric.Name,
			MType: typeMetricName,
		}

		err = h.agentService.SendMetric(MetricGaugeUpdateRequest)
		if err != nil {
			return fmt.Errorf("failed to send metric %s to server: %w", metric.Name, err)
		}
	}

	return nil
}
