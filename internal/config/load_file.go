package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadConfigFromFile[T any](path string, cfg *T) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}
