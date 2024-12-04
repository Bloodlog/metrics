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

func Handle(configs *config.Config, repository *repository.Repository) error {
	serverAddr := "http://" + fmt.Sprintf("%s:%d", configs.NetAddress.Host, configs.NetAddress.Port)
	client := resty.New().SetBaseURL(serverAddr)
	metricsChan := make(chan []string)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		getMetrics(metricsChan, repository, configs.PollInterval)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendMetrics(metricsChan, client, configs.ReportInterval, configs.Debug)
	}()

	wg.Wait()

	return nil
}

func getMetrics(metricsChan chan []string, repository *repository.Repository, pollInterval int) {
	for {
		metrics := repository.GetMemoryMetrics()

		var stringMetrics []string
		for _, metric := range metrics {
			stringMetrics = append(stringMetrics, fmt.Sprintf("%s:%d", metric.Name, metric.Value))
		}

		metricsChan <- stringMetrics
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendMetrics(metricsChan chan []string, client *resty.Client, reportIntervalrepository int, debug bool) {
	counter := 0
	for {
		time.Sleep(time.Duration(reportIntervalrepository) * time.Second)
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
				fmt.Println("invalid metric string format:", err)
				return
			}
			err2 := service.SendMetric(client, parts[0], parts[1], debug)
			if err2 != nil {
				fmt.Println("Error sending metrics:", err2)
				return
			}
		}
	}
}
