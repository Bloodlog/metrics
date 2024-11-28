package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestCounterService_SendGauge(t *testing.T) {
	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())
	metricName := "metricName"
	metricValue := uint64(5)
	url := "http://localhost:8080/update/gauge/" + metricName + "/" + strconv.Itoa(int(metricValue))
	responder := httpmock.NewStringResponder(http.StatusOK, ``)
	httpmock.RegisterResponder("POST", url, responder)

	err := SendMetric(client, metricName, metricValue, false)

	assert.NoError(t, err)
}
