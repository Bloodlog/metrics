package handlers

import (
	"metrics/internal/agent/clients"
	"metrics/internal/agent/config"
	"metrics/internal/agent/repository"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func setupTestClient() *clients.Client {
	logger := zap.NewNop().Sugar()
	clientResty := resty.New()
	httpmock.ActivateNonDefault(clientResty.GetClient())

	return &clients.Client{
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

	h := NewHandlers(
		client,
		&config.Config{},
		repository.NewMemoryRepository(),
		repository.NewSystemRepository(),
		client.Logger,
	)
	err := h.sendAPI([]repository.Metric{{Name: "metric1", Value: 10}}, 5)
	assert.NoError(t, err)
}

func TestSendBatch_Success(t *testing.T) {
	client := setupTestClient()
	defer httpmock.DeactivateAndReset()

	url := "/updates"
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewStringResponder(http.StatusOK, ""))

	h := NewHandlers(
		client,
		&config.Config{},
		repository.NewMemoryRepository(),
		repository.NewSystemRepository(),
		client.Logger,
	)

	metrics := []repository.Metric{
		{Name: "metric1", Value: 10},
		{Name: "metric2", Value: 20},
	}

	err := h.sendBatch(metrics, 5)
	assert.NoError(t, err)
}
