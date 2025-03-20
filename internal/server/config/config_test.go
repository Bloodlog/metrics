package config

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
		true,
	)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.NetAddress.Host)
	assert.Equal(t, "8080", cfg.NetAddress.Port)

	assert.Equal(t, 500, cfg.StoreInterval)

	assert.Equal(t, true, cfg.Restore)

	assert.Equal(t, "metrics_temp.json", cfg.FileStoragePath)

	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.DatabaseDsn)
	assert.Equal(t, "my-secret-key", cfg.Key)

	assert.Equal(t, true, cfg.Debug)
}
