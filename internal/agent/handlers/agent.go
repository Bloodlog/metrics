package handlers

import (
	"github.com/go-resty/resty/v2"
	"metrics/internal/agent/repository"
	"metrics/internal/agent/service"
	"time"
)

func Handle(repository *repository.Repository, debug bool) error {
	sleepTimeSec := 2
	counter := 0

	client := resty.New()

	for {
		time.Sleep(time.Duration(sleepTimeSec) * time.Second)
		metrics := repository.GetMemoryMetrics()
		counter++
		err := service.SendIncrement(client, uint64(counter), debug)
		if err != nil {
			return err
		}

		err = sendMetrics(client, metrics, debug)
		if err != nil {
			return err
		}
	}
}

func sendMetrics(client *resty.Client, metrics []repository.Metric, debug bool) error {
	for _, metric := range metrics {
		err := service.SendMetric(client, metric.Name, metric.Value, debug)
		if err != nil {
			return err
		}
	}

	return nil
}
