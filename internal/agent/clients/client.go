package clients

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	handlerLogger := logger.With("client", "send request")
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	return resty.New().
		SetBaseURL(serverAddr).
		SetRetryCount(len(retryIntervals)).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= http.StatusInternalServerError
		}).
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			return handleRetry(req, retryIntervals, handlerLogger)
		}).
		OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
			handlerLogger.Infof("Received response from %s with status: %d, body: %v",
				resp.Request.URL, resp.StatusCode(), resp.String())
			return nil
		}).
		OnError(func(req *resty.Request, err error) {
			handlerLogger.Infoln("Request to %s failed: %v", req.URL, err)
		})
}

func handleRetry(
	req *resty.Request,
	retryIntervals []time.Duration,
	logger *zap.SugaredLogger,
) error {
	attempt := req.Attempt - 1
	if attempt > 0 && attempt <= len(retryIntervals) {
		logger.Infof("Retrying request to %s (attempt %d/%d), waiting %s",
			req.URL, attempt, len(retryIntervals), retryIntervals[attempt-1])
		time.Sleep(retryIntervals[attempt-1])
	}
	return nil
}
