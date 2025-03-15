package config

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseAddressFlags(t *testing.T) {
	resetFlags()
	flagAddress := "a"
	host := "localhost"
	port := "8080"
	os.Args = []string{"cmd", "--" + flagAddress + "=" + host + ":" + port}
	cfg, _ := ParseFlags()
	assert.Equal(t, host, cfg.NetAddress.Host)
	assert.Equal(t, port, cfg.NetAddress.Port)
}

func TestParseIntervalFlags(t *testing.T) {
	resetFlags()
	flagInterval := "i"
	storeInterval := "500"
	os.Args = []string{"cmd", "--" + flagInterval + "=" + storeInterval}
	cfg, _ := ParseFlags()
	storeIntervalInt, _ := strconv.Atoi(storeInterval)

	assert.Equal(t, storeIntervalInt, cfg.StoreInterval)
}

func TestParseRestoreFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--r=true"}
	cfg, _ := ParseFlags()

	assert.Equal(t, true, cfg.Restore)
}

func TestParseStoragePathFlags(t *testing.T) {
	resetFlags()
	flagStoragePath := "f"
	storagePath := "metrics_temp.json"
	os.Args = []string{"cmd", "--" + flagStoragePath + "=" + storagePath}
	cfg, _ := ParseFlags()

	assert.Equal(t, storagePath, cfg.FileStoragePath)
}

func TestParseDatabaseFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--d=postgres://user:pass@localhost:5432/dbname"}
	cfg, _ := ParseFlags()

	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.DatabaseDsn)
}

func TestParseKeyFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--k=my-secret-key"}
	cfg, _ := ParseFlags()

	assert.Equal(t, "my-secret-key", cfg.Key)
}

func TestParsePprofFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--pprof=true"}
	cfg, _ := ParseFlags()

	assert.Equal(t, true, cfg.Debug)
}
