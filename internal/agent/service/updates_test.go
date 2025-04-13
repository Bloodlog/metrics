package service

import (
	"metrics/internal/agent/dto"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSendMetricsBatch(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())

	metricsUpdateRequest := dto.MetricsUpdateRequests{
		Metrics: []dto.MetricsUpdateRequest{
			{
				ID:    "metric1",
				MType: "gauge",
				Value: new(float64),
			},
		},
	}

	metricsUpdateRequest.Metrics[0].Value = new(float64)
	*metricsUpdateRequest.Metrics[0].Value = 10.5

	url := "/updates"
	responder := httpmock.NewStringResponder(http.StatusOK, "")
	httpmock.RegisterResponder("POST", url, responder)

	err := SendMetricsBatch(client, metricsUpdateRequest)

	assert.NoError(t, err)

	reqs := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, reqs["POST "+url])
}
