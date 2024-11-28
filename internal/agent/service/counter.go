package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func SendIncrement(client *resty.Client, counter uint64, debug bool) error {
	response, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"counter": strconv.Itoa(int(counter)),
		}).
		Post("http://localhost:8080/update/counter/PollCount/{counter}")
	if err != nil {
		return err
	}

	if response.IsError() {
		return err
	}

	if debug {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		log := "[" + timeStr + "] " + response.Request.URL + " " + strconv.Itoa(response.StatusCode())
		fmt.Println(log)
	}

	return nil
}
