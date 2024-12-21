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
	var metricCounterRequest MetricsCounterRequest
	delta := int64(counter)

	metricCounterRequest = MetricsCounterRequest{
		Delta: &delta,
		ID:    "PoolCounter",
		MType: "counter",
	}

	err := SendIncrement(client, metricCounterRequest)

	assert.NoError(t, err)
}
