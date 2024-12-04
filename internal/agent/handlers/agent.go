package handlers

import (
	"fmt"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

func Handle(configs *config.Config, storage *repository.Repository) error {
	serverAddr := "http://" + fmt.Sprintf("%s:%d", configs.NetAddress.Host, configs.NetAddress.Port)
	client := resty.New().SetBaseURL(serverAddr)
	metricsChan := make(chan []string)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		getMetrics(metricsChan, storage, configs.PollInterval)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := sendMetrics(metricsChan, client, configs.ReportInterval, configs.Debug)
		if err != nil {
			return
		}
	}()

	wg.Wait()

	return nil
}

func getMetrics(metricsChan chan []string, storage *repository.Repository, pollInterval int) {
	for {
		metrics := storage.GetMemoryMetrics()

		var stringMetrics []string
		for _, metric := range metrics {
			stringMetrics = append(stringMetrics, fmt.Sprintf("%s:%d", metric.Name, metric.Value))
		}

		metricsChan <- stringMetrics
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendMetrics(metricsChan chan []string, client *resty.Client, reportIntervalrepository int, debug bool) error {
	const numberParts = 2
	counter := 0
	for {
		time.Sleep(time.Duration(reportIntervalrepository) * time.Second)
		counter++
		err := service.SendIncrement(client, uint64(counter), debug)
		if err != nil {
			return fmt.Errorf("error sending increment: %w", err)
		}
		stringMetrics := <-metricsChan
		for _, metricStr := range stringMetrics {
			parts := strings.Split(metricStr, ":")
			if len(parts) != numberParts {
				return fmt.Errorf("invalid metric string format: %w", err)
			}
			err := service.SendMetric(client, parts[0], parts[1], debug)
			if err != nil {
				return fmt.Errorf("error sending metrics: %w", err)
			}
		}
	}
}
