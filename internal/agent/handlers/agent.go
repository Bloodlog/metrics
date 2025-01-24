package handlers

import (
	"fmt"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

const (
	maxNumberAttempts  = 3
	retryWaitSecond    = 2
	retryMaxWaitSecond = 5
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
	client := createClient(serverAddr, logger)

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
			var metricCounterRequest service.MetricsCounterRequest
			delta := int64(counter)

			metricCounterRequest = service.MetricsCounterRequest{
				Delta: &delta,
				ID:    nameCounter,
				MType: typeCounter,
			}

			err := service.SendIncrement(h.client, metricCounterRequest)

			counter = 0
			if err != nil {
				return fmt.Errorf("failed to send Increment %d to server: %w", counter, err)
			}

			for _, metric := range metrics {
				var MetricGaugeUpdateRequest service.MetricsUpdateRequest
				valueFloat := float64(metric.Value)

				MetricGaugeUpdateRequest = service.MetricsUpdateRequest{
					Value: &valueFloat,
					ID:    metric.Name,
					MType: typeMetricName,
				}

				err := service.SendMetric(h.client, MetricGaugeUpdateRequest)
				if err != nil {
					return fmt.Errorf("failed to send metric %s to server: %w", metric.Name, err)
				}
			}
		}
	}
}

func createClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	handlerLogger := logger.With("client", "send request")
	return resty.New().
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetRetryCount(maxNumberAttempts).
		SetRetryWaitTime(retryWaitSecond * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitSecond * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= http.StatusInternalServerError
		}).
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			handlerLogger.Infof("Sending request to %s with body: %v", req.URL, req.Body)
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
			handlerLogger.Infof("Received response from %s with status: %d, body: %v",
				resp.Request.URL, resp.StatusCode(), resp.String())
			return nil
		}).
		OnError(func(req *resty.Request, err error) {
			handlerLogger.Infoln("Request to %s failed: %v", req.URL, err)
		})
}
