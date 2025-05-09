package handlers

import (
	"metrics/internal/config"
	repository2 "metrics/internal/repository"
	"metrics/internal/service"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func setupTestClient() *service.Client {
	logger := zap.NewNop().Sugar()
	clientResty := resty.New()
	httpmock.ActivateNonDefault(clientResty.GetClient())

	return &service.Client{
		RestyClient: clientResty,
		Logger:      logger,
		Key:         "key",
	}
}

func TestSendAPI_Success(t *testing.T) {
	client := setupTestClient()
	defer httpmock.DeactivateAndReset()

	url := "/update/"
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewStringResponder(http.StatusOK, ""))

	h := NewAgentHandler(
		client,
		&config.AgentConfig{},
		repository2.NewMemoryRepository(),
		repository2.NewSystemRepository(),
		client.Logger,
	)
	err := h.sendAPI([]repository2.Metric{{Name: "metric1", Value: 10}}, 5)
	assert.NoError(t, err)
}

func TestSendBatch_Success(t *testing.T) {
	client := setupTestClient()
	defer httpmock.DeactivateAndReset()

	url := "/updates"
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewStringResponder(http.StatusOK, ""))

	h := NewAgentHandler(
		client,
		&config.AgentConfig{},
		repository2.NewMemoryRepository(),
		repository2.NewSystemRepository(),
		client.Logger,
	)

	metrics := []repository2.Metric{
		{Name: "metric1", Value: 10},
		{Name: "metric2", Value: 20},
	}

	err := h.sendBatch(metrics, 5)
	assert.NoError(t, err)
}
