package handlers

import (
	"fmt"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

const maxNumberAttempts = 3
const retryWaitSecond = 2
const retryMaxWaitSecond = 5

func Handle(configs *config.Config, storage *repository.Repository, logger zap.SugaredLogger) error {
	serverAddr := "http://" + net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)
	client := resty.New().
		SetBaseURL(serverAddr).
		SetRetryCount(maxNumberAttempts).
		SetRetryWaitTime(retryWaitSecond * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitSecond * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= 500
		}).
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			logger.Infof("Sending request to %s with body: %v", req.URL, req.Body)
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
			logger.Infof("Received response from %s with status: %d, body: %v",
				resp.Request.URL, resp.StatusCode(), resp.String())
			return nil
		}).
		OnError(func(req *resty.Request, err error) {
			logger.Infoln("Request to %s failed: %v", req.URL, err)
		})

	pollTicker := time.NewTicker(time.Duration(configs.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(configs.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	var metrics []repository.Metric
	counter := 0

	for {
		select {
		case <-pollTicker.C:
			metrics = storage.GetMemoryMetrics()
			counter++

		case <-reportTicker.C:
			var metricCounterRequest service.MetricsCounterRequest
			delta := int64(counter)

			metricCounterRequest = service.MetricsCounterRequest{
				Delta: &delta,
				ID:    "PollCount",
				MType: "counter",
			}

			err := service.SendIncrement(client, metricCounterRequest)

			counter = 0
			if err != nil {
				logger.Infoln(err.Error(), "handler", "send Increment")
				return fmt.Errorf("failed to send Increment %d to server: %w", counter, err)
			}

			for _, metric := range metrics {
				var MetricGaugeUpdateRequest service.MetricsUpdateRequest
				valueFloat := float64(metric.Value)

				MetricGaugeUpdateRequest = service.MetricsUpdateRequest{
					Value: &valueFloat,
					ID:    metric.Name,
					MType: "gauge",
				}

				err := service.SendMetric(client, MetricGaugeUpdateRequest)
				if err != nil {
					logger.Infoln(err.Error(), "handler", "send metric")
					return fmt.Errorf("failed to send metric %s to server: %w", metric.Name, err)
				}
			}
		}
	}
}
