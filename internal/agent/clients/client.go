package clients

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	handlerLogger := logger.With("client", "send request")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = CustomBackoff

	client := resty.New().
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	clientRetry := client.
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			handlerLogger.Info("retry request")
			httpReq, err := retryablehttp.NewRequest(req.Method, req.URL, req.Body)
			if err != nil {
				return fmt.Errorf("failed to create retryable request: %w", err)
			}

			resp, err := retryClient.Do(httpReq)
			if err != nil {
				return fmt.Errorf("failed to send request with retries: %w", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					handlerLogger.Infoln("failed to close response body: %v", err)
				}
			}()

			return nil
		})

	return clientRetry
}

func CustomBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	backoffIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if attemptNum-1 < len(backoffIntervals) {
		return backoffIntervals[attemptNum-1]
	}
	return max
}
