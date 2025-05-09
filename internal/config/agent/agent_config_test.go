package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAgentFlags(t *testing.T) {
	cfg, err := processAgentFlags(
		"localhost:8080",
		500,
		600,
		"my-secret-key",
		10,
		"test",
		"",
		"",
	)
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost:8080", cfg.Address)

	assert.Equal(t, 500, cfg.ReportInterval)
	assert.Equal(t, 600, cfg.PollInterval)

	assert.Equal(t, "my-secret-key", cfg.Key)
	assert.Equal(t, 10, cfg.RateLimit)
	assert.Equal(t, "test", cfg.CryptoKey)
}

func TestParseAgentFlags(t *testing.T) {
	cfg, err := ParseAgentFlags()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080", cfg.Address)

	assert.Equal(t, 10, cfg.ReportInterval)
	assert.Equal(t, 2, cfg.PollInterval)

	assert.Equal(t, "", cfg.Key)
	assert.Equal(t, 1, cfg.RateLimit)
	assert.Equal(t, "", cfg.CryptoKey)
}
