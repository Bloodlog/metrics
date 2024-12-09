package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultAddress         = "http://localhost:8080"
	EnvAddress             = "ADDRESS"
	AddressFlagDescription = "HTTP server address in the format host:port (default: localhost:8080)"
)

type Config struct {
	NetAddress NetAddress
	Debug      bool
}

type NetAddress struct {
	Host string
	Port string
}

func ParseFlags() (*Config, error) {
	addressFlag := flag.String("a", DefaultAddress, AddressFlagDescription)
	flag.Parse()

	uknownArguments := flag.Args()
	if err := validateUnknownArgs(uknownArguments); err != nil {
		log.Printf("error: unknown flags or arguments detected: %v", uknownArguments)
		return nil, err
	}

	finalAddress, err := getStringValue(*addressFlag, EnvAddress)
	if err != nil {
		log.Printf("error: invalid address: %v", err)
		return nil, err
	}

	host, port, err := parseAddress(finalAddress)
	if err != nil {
		log.Printf("error: invalid address: %v", err)
		return nil, err
	}

	return &Config{
		NetAddress: NetAddress{Host: host, Port: port},
		Debug:      false,
	}, nil
}

func validateUnknownArgs(unknownArgs []string) error {
	if len(unknownArgs) > 0 {
		return errors.New("unknown flags or arguments detected")
	}
	return nil
}

func parseAddress(address string) (string, string, error) {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", "", errors.New("failed to parse address (expected host:port)")
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	if host == "" || port == "" {
		return "", "", errors.New("failed to parse address (expected host:port)")
	}

	return host, port, nil
}

func getStringValue(flagValue, envKey string) (string, error) {
	if envValue, exists := os.LookupEnv(envKey); exists {
		return envValue, nil
	}

	if flagValue != "" {
		return flagValue, nil
	}

	return "", fmt.Errorf("missing required configuration: %s or flag value", envKey)
}
