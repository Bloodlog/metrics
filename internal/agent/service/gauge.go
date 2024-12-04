package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func SendMetric(client *resty.Client, name string, value string, debug bool) error {
	response, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"metricName":  name,
			"metricValue": value,
		}).
		Post("/update/gauge/{metricName}/{metricValue}")
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}

	if debug {
		timeStr := time.Now().Format(time.DateTime)
		log := "[" + timeStr + "] " + response.Request.URL + " " + strconv.Itoa(response.StatusCode())
		fmt.Println(log)
	}

	return nil
}
