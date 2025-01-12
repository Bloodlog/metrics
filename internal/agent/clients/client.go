package clients

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, logger *zap.SugaredLogger) *resty.Client {
	const (
		maxNumberAttempts  = 3
		retryWaitSecond    = 2
		retryMaxWaitSecond = 5
	)
	handlerLogger := logger.With("client", "send request")
	return resty.New().
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetRetryCount(maxNumberAttempts).
		SetRetryWaitTime(retryWaitSecond * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitSecond * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= http.StatusInternalServerError
		}).
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			handlerLogger.Infof("Sending request to %s with body: %v", req.URL, req.Body)
			return nil
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
