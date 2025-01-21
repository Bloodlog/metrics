package clients

import (
	"context"
	"errors"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	handlerLogger := logger.With("agent", "client")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = customBackoff
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
				handlerLogger.Infoln(
					"status_code", resp.StatusCode,
					"retryable", true,
				)
				return true, nil
			}

			var DNSError *net.DNSError
			retry := errors.Is(err, syscall.ECONNREFUSED) ||
				errors.Is(err, syscall.ETIMEDOUT) ||
				err.Error() == "EOF" ||
				errors.As(err, &DNSError)

			if retry {
				handlerLogger.Infoln("Connect problem", "retryable", true)
				return true, nil
			}

			if resp == nil {
				handlerLogger.Infoln("Connect problem and response == nil and EOF", "retryable", true)
				return true, nil
			}
		}

		handlerLogger.Infoln("Client", "CheckRetry", "retryable", false)
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
