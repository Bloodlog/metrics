package clients

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = customBackoff
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ETIMEDOUT) {
				return true, nil
			}
		}
		return false, nil
	}

	restyClient := resty.NewWithClient(retryClient.StandardClient()).
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	return restyClient
}

func customBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	backoffIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if attemptNum-1 >= 0 && attemptNum-1 < len(backoffIntervals) {
		return backoffIntervals[attemptNum-1]
	}
	return max
}
