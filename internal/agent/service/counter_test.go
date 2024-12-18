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

	responder := httpmock.NewStringResponder(http.StatusOK, ``)
	url := "/update/counter/PollCount/42"

	httpmock.RegisterResponder("POST", url, responder)

	err := SendIncrement(client, 42)

	assert.NoError(t, err)
}
