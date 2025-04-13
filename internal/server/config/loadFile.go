package config

import (
	"encoding/json"
	"fmt"
	"io"
	"metrics/internal/server/dto"
	"os"
)

func loadConfigFromFile(path string) (*dto.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg dto.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
