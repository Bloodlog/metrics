package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")

	if err := os.WriteFile(tmpFile, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return tmpFile
}

func TestLoadAgentConfigFromFile(t *testing.T) {
	jsonData := `{
		"key": "secret",
		"crypto_key": "rsa_key",
		"address": "127.0.0.1:8080",
		"report_interval": 10,
		"poll_interval": 2,
		"rate_limit": 100,
		"batch": true
	}`

	path := writeTempFile(t, jsonData)

	var cfg AgentConfig
	err := LoadConfigFromFile(path, &cfg)
	if err != nil {
		t.Fatalf("LoadAgentConfigFromFile failed: %v", err)
	}

	if cfg.Address != "127.0.0.1:8080" {
		t.Errorf("expected Address to be '127.0.0.1:8080', got %q", cfg.Address)
	}
}

func TestLoadServerConfigFromFile(t *testing.T) {
	jsonData := `{
		"key": "serverkey",
		"address": "0.0.0.0:9000",
		"store_file": "/tmp/data.json",
		"database_dsn": "postgres://user:pass@localhost/db",
		"crypto_key": "server_rsa",
		"store_interval": 30,
		"restore": true,
		"debug": true
	}`

	path := writeTempFile(t, jsonData)

	var cfg ServerConfig
	err := LoadConfigFromFile(path, &cfg)
	if err != nil {
		t.Fatalf("LoadServerConfigFromFile failed: %v", err)
	}

	if !cfg.Restore {
		t.Errorf("expected Restore to be true")
	}

	if cfg.StoreInterval != 30 {
		t.Errorf("expected StoreInterval to be 30, got %d", cfg.StoreInterval)
	}
}
