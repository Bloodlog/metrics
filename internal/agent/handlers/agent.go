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
			metricsRequests := service.MetricsUpdateRequests{}
			delta := int64(counter)

			metric := service.MetricsUpdateRequest{
				Delta: &delta,
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
			counter = 0
		}
	}
}
