package handlers

import (
	"fmt"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"net"
	"strconv"
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
			err := service.SendIncrement(client, uint64(counter))
			if err != nil {
				return fmt.Errorf("failed to send increment: %w", err)
			}

			for _, metric := range metrics {
				metricValueString := strconv.FormatUint(metric.Value, 10)
				err := service.SendMetric(client, metric.Name, metricValueString)
				if err != nil {
					return fmt.Errorf("failed to send metric %s: %w", metric.Name, err)
				}
			}
		}
	}
}
