package repository

import (
	"context"
	"metrics/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewMetricStorage(t *testing.T) {
	logger := zap.NewExample().Sugar()

	tests := []struct {
		cfg  *config.ServerConfig
		name string
	}{
		{
			name: "Memory storage",
			cfg:  &config.ServerConfig{},
		},
		{
			name: "File retry storage",
			cfg: &config.ServerConfig{
				FileStoragePath: "/some/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewMetricStorage(context.Background(), tt.cfg, logger)

			assert.NoError(t, err)
			assert.NotNil(t, storage)
		})
	}
}
