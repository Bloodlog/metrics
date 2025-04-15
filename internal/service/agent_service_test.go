package service

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCounterService_SendIncrement(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())

	responder := httpmock.NewStringResponder(http.StatusOK, "")
	url := "/update/"

	httpmock.RegisterResponder("POST", url, responder)

	counter := 42
	var metricCounterRequest AgentMetricsCounterRequest
	delta := int64(counter)

	metricCounterRequest = AgentMetricsCounterRequest{
		Delta: &delta,
		ID:    "PoolCounter",
		MType: "counter",
	}

	err := SendIncrement(client, metricCounterRequest)

	assert.NoError(t, err)
}

func TestSendMetricsBatch(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())

	metricsUpdateRequest := AgentMetricsUpdateRequests{
		Metrics: []AgentMetricsUpdateRequest{
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

func TestCounterService_SendGauge(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())
	metricName := "metricName"
	metricValue := 5
	url := "/update/"
	responder := httpmock.NewStringResponder(http.StatusOK, "")
	httpmock.RegisterResponder("POST", url, responder)

	var MetricGaugeUpdateRequest AgentMetricsGaugeUpdateRequest
	valueFloat := float64(metricValue)

	MetricGaugeUpdateRequest = AgentMetricsGaugeUpdateRequest{
		Value: &valueFloat,
		ID:    metricName,
		MType: "gauge",
	}

	err := SendMetric(client, MetricGaugeUpdateRequest)

	assert.NoError(t, err)
}
