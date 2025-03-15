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

func TestParseReportIntervalFlags(t *testing.T) {
	resetFlags()
	flagInterval := "r"
	interval := "500"
	os.Args = []string{"cmd", "--" + flagInterval + "=" + interval}
	cfg, _ := ParseFlags()
	intervalInt, _ := strconv.Atoi(interval)

	assert.Equal(t, intervalInt, cfg.ReportInterval)
}

func TestParsePoolIntervalFlags(t *testing.T) {
	resetFlags()
	flagInterval := "p"
	interval := "600"
	os.Args = []string{"cmd", "--" + flagInterval + "=" + interval}
	cfg, _ := ParseFlags()
	intervalInt, _ := strconv.Atoi(interval)

	assert.Equal(t, intervalInt, cfg.PollInterval)
}

func TestParseKeyFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--k=my-secret-key"}
	cfg, _ := ParseFlags()

	assert.Equal(t, "my-secret-key", cfg.Key)
}
