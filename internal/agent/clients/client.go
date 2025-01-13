package clients

import (
	"metrics/internal/agent/config"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func CreateClient(serverAddr string, configs config.ClientSetting, logger *zap.SugaredLogger) *resty.Client {
	handlerLogger := logger.With("client", "send request")
	return resty.New().
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetRetryCount(configs.MaxNumberAttempts).
		SetRetryWaitTime(time.Duration(configs.RetryWaitSecond) * time.Second).
		SetRetryMaxWaitTime(time.Duration(configs.RetryMaxWaitSecond) * time.Second).
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
