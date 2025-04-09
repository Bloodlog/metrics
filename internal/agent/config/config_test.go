package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFlags(t *testing.T) {
	cfg, err := processFlags("localhost:8080", 500, 600, "my-secret-key", 10)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.NetAddress.Host)
	assert.Equal(t, "8080", cfg.NetAddress.Port)

	assert.Equal(t, 500, cfg.ReportInterval)
	assert.Equal(t, 600, cfg.PollInterval)

	assert.Equal(t, "my-secret-key", cfg.Key)
	assert.Equal(t, 10, cfg.RateLimit)
}

func TestParseFlags(t *testing.T) {
	cfg, err := ParseFlags()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.NetAddress.Host)
	assert.Equal(t, "8080", cfg.NetAddress.Port)

	assert.Equal(t, 10, cfg.ReportInterval)
	assert.Equal(t, 2, cfg.PollInterval)

	assert.Equal(t, "", cfg.Key)
	assert.Equal(t, 1, cfg.RateLimit)
}
