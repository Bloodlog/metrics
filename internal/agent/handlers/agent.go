package handlers

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"strings"
	"time"
)

func Handle(serverAddr string, reportIntervalrepository int, pollInterval int, repository *repository.Repository, debug bool) error {
	metricsChan := make(chan []string)
	go func() {
		for {
			time.Sleep(time.Duration(reportIntervalrepository) * time.Second)
			metrics := repository.GetMemoryMetrics()
			var stringMetrics []string
			for _, metric := range metrics {
				stringMetrics = append(stringMetrics, fmt.Sprintf("%s:%d", metric.Name, metric.Value))
			}

			metricsChan <- stringMetrics
		}
	}()

	go func() {
		counter := 0
		client := resty.New().SetBaseURL(serverAddr)
		for {
			time.Sleep(time.Duration(pollInterval) * time.Second)
			counter++
			err := service.SendIncrement(client, uint64(counter), debug)
			if err != nil {
				fmt.Println("Error sending increment:", err)
				return
			}
			stringMetrics := <-metricsChan
			for _, metricStr := range stringMetrics {
				parts := strings.Split(metricStr, ":")
				if len(parts) != 2 {
					fmt.Println("неверный формат строки метрики:", err)
					return
				}
				err2 := service.SendMetric(client, parts[0], parts[1], debug)
				if err != nil {
					fmt.Println("Error sending metrics:", err2)
					return
				}
			}
		}
	}()

	return nil
}
