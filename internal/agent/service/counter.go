package service

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func SendIncrement(client *resty.Client, counter uint64) error {
	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"counter": strconv.Itoa(int(counter)),
		}).
		Post("/update/counter/PollCount/{counter}")
	if err != nil {
		return fmt.Errorf("failed to send increment: %w", err)
	}

	return nil
}
