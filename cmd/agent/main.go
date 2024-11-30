package main

import (
	"flag"
	"fmt"
	"log"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type Config struct {
	NetAddress     NetAddress
	ReportInterval int
	PollInterval   int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := parseFlags()
	if err != nil {
		return err
	}

	rep := repository.NewRepository()
	debug := false
	schema := "http://"

	serverAddr := schema + fmt.Sprintf("%s:%d", config.NetAddress.Host, config.NetAddress.Port)
	return handlers.Handle(serverAddr, config.ReportInterval, config.PollInterval, rep, debug)
}

func parseFlags() (*Config, error) {
	address := flag.String("a", "localhost:8080", "Адрес HTTP-сервера в формате host:port (по умолчанию localhost:8080)")
	reportInterval := flag.Int("r", 10, "Позволяет переопределять частоту отправки метрик на сервер (по умолчанию 10 секунд)")
	pollInterval := flag.Int("p", 2, "Позволяет переопределять частоту опроса метрик из пакета runtime (по умолчанию 2 секунды)")

	flag.Parse()

	if len(flag.Args()) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Ошибка: обнаружены неизвестные флаги")
		os.Exit(1)
	}

	parts := strings.Split(*address, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("неверный формат адреса: %s (ожидается host:port)", *address)
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать порт в число: %w", err)
	}

	return &Config{
		NetAddress:     NetAddress{Host: host, Port: port},
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
	}, nil
}
