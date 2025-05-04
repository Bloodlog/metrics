package handlers

import (
	"context"
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

	agentService := service.NewHTTPMetricSender(client.RestyClient)
	h := NewAgentHandler(
		&config.AgentConfig{},
		repository2.NewMemoryRepository(),
		repository2.NewSystemRepository(),
		agentService,
		client.Logger,
	)
	ctx := context.Background()
	err := h.sendAPI(ctx, []repository2.Metric{{Name: "metric1", Value: 10}}, 5)
	assert.NoError(t, err)
}

func TestSendBatch_Success(t *testing.T) {
	client := setupTestClient()
	defer httpmock.DeactivateAndReset()

	url := "/updates"
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewStringResponder(http.StatusOK, ""))

	agentService := service.NewHTTPMetricSender(client.RestyClient)
	h := NewAgentHandler(
		&config.AgentConfig{},
		repository2.NewMemoryRepository(),
		repository2.NewSystemRepository(),
		agentService,
		client.Logger,
	)

	metrics := []repository2.Metric{
		{Name: "metric1", Value: 10},
		{Name: "metric2", Value: 20},
	}

	ctx := context.Background()
	err := h.sendBatch(ctx, metrics, 5)
	assert.NoError(t, err)
}
