package service

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCounterService_SendGauge(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())
	metricName := "metricName"
	metricValue := 5
	url := "/update"
	responder := httpmock.NewStringResponder(http.StatusOK, "")
	httpmock.RegisterResponder("POST", url, responder)

	var MetricGaugeUpdateRequest MetricsUpdateRequest
	valueFloat := float64(metricValue)

	MetricGaugeUpdateRequest = MetricsUpdateRequest{
		Value: &valueFloat,
		ID:    metricName,
		MType: "gauge",
	}

	err := SendMetric(client, MetricGaugeUpdateRequest)

	assert.NoError(t, err)
}
