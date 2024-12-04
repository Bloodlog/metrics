package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func SendIncrement(client *resty.Client, counter uint64, debug bool) error {
	response, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"counter": strconv.Itoa(int(counter)),
		}).
		Post("/update/counter/PollCount/{counter}")
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
