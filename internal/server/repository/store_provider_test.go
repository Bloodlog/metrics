package repository

import (
	"context"
	"testing"

	"metrics/internal/server/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewMetricStorage(t *testing.T) {
	logger := zap.NewExample().Sugar()

	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "Memory storage",
			cfg:  &config.Config{},
		},
		{
			name: "File retry storage",
			cfg: &config.Config{
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
