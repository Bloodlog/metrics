package service

import (
	"errors"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

var (
	ErrSendIncrement = errors.New("failed to send POST request PollCount")
)

func SendIncrement(client *resty.Client, counter uint64) error {
	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"counter": strconv.Itoa(int(counter)),
		}).
		Post("/update/counter/PollCount/{counter}")
	if err != nil {
		log.Printf("failed to send POST request PollCount: %v", err)
		return ErrSendIncrement
	}

	return nil
}
