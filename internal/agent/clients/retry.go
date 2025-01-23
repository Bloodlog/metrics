package clients

import (
	"context"
	"errors"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

func createRetryPolicy(logger *zap.SugaredLogger) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
				logger.Infoln(
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
				logger.Infoln("Connect problem", "retryable", true)
				return true, nil
			}

			if resp == nil {
				logger.Infoln("Connect problem and response == nil and EOF", "retryable", true)
				return true, nil
			}
		}

		return false, nil
	}
}

func customBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	backoffIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if attemptNum-1 >= 0 && attemptNum-1 < len(backoffIntervals) {
		return backoffIntervals[attemptNum-1]
	}
	return max
}
