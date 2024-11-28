package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCounterService_SendIncrement(t *testing.T) {

	client := resty.New()

	httpmock.ActivateNonDefault(client.GetClient())

	responder := httpmock.NewStringResponder(http.StatusOK, ``)
	url := "http://localhost:8080/update/counter/PollCount/42"

	httpmock.RegisterResponder("POST", url, responder)

	err := SendIncrement(client, 42, false)

	assert.NoError(t, err)
}
