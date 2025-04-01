package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	serverAddr := "http://example.com"
	key := "test-key"

	client := NewClient(serverAddr, key, sugar)

	assert.NotNil(t, client)

	assert.NotNil(t, client.RestyClient)

	assert.Equal(t, sugar, client.Logger)

	assert.Equal(t, key, client.Key)

	header := client.RestyClient.Header
	assert.Equal(t, serverAddr, client.RestyClient.BaseURL)
	assert.Contains(t, header.Get("Content-Encoding"), "gzip")
	assert.Contains(t, header.Get("Content-Type"), "application/json")
}
