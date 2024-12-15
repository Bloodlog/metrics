package handlers

import (
	"fmt"
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
	"time"

	"github.com/go-resty/resty/v2"
)

func Handle(configs *config.Config, storage *repository.Repository) error {
	serverAddr := "http://" + net.JoinHostPort(configs.NetAddress.Host, configs.NetAddress.Port)
	client := resty.New().SetBaseURL(serverAddr)

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
				ID:    "PoolCounter",
				MType: "counter",
			}

			err := service.SendIncrement(client, metricCounterRequest)

			counter = 0
			if err != nil {
				log.Printf("failed to send POST request Increment: %v", err)
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
					log.Printf("failed to send POST metric: %v", err)
					return fmt.Errorf("failed to send metric %s to server: %w", metric.Name, err)
				}
			}
		}
	}
}
