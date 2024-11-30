package main

import (
	"flag"
	"fmt"
	"log"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type NetAddress struct {
	Host string
	Port int
}

func run() error {
	netAddr, err := parseFlags()
	if err != nil {
		return err
	}
	memStorage := repository.NewMemStorage()
	debug := false

	serverAddr := fmt.Sprintf("%s:%d", netAddr.Host, netAddr.Port)
	fmt.Println("Запуск сервера на адресе:", serverAddr)
	return router.Run(serverAddr, memStorage, debug)
}

func parseFlags() (*NetAddress, error) {
	address := flag.String("a", "localhost:8080", "Адрес HTTP-сервера в формате host:port (по умолчанию localhost:8080)")

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

	return &NetAddress{Host: host, Port: port}, nil
}
