package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddressFlags(t *testing.T) {
	cfg, err := processFlags(
		"localhost:8080",
		500,
		"metrics_temp.json",
		true,
		"postgres://user:pass@localhost:5432/dbname",
		"my-secret-key",
		"test",
		true,
		"",
		"",
		"",
	)
	assert.NoError(t, err)

	assert.Equal(t, "localhost:8080", cfg.Address)

	assert.Equal(t, 500, cfg.StoreInterval)

	assert.Equal(t, true, cfg.Restore)

	assert.Equal(t, "metrics_temp.json", cfg.FileStoragePath)

	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.DatabaseDsn)
	assert.Equal(t, "my-secret-key", cfg.Key)

	assert.Equal(t, "test", cfg.CryptoKey)

	assert.Equal(t, true, cfg.Debug)
}

func TestParseFlags(t *testing.T) {
	cfg, err := ParseFlags()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", cfg.Address)

	assert.Equal(t, 300, cfg.StoreInterval)

	assert.Equal(t, true, cfg.Restore)

	assert.Equal(t, "metrics.json", cfg.FileStoragePath)

	assert.Equal(t, "", cfg.DatabaseDsn)
	assert.Equal(t, "", cfg.Key)

	assert.Equal(t, "", cfg.CryptoKey)

	assert.Equal(t, false, cfg.Debug)
}
