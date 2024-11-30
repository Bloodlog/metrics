package main

import (
	"fmt"
	"log"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := parseFlag()
	if err != nil {
		return err
	}

	rep := repository.NewRepository()
	debug := true

	serverAddr := fmt.Sprintf("%s:%d", config.NetAddress.Host, config.NetAddress.Port)
	return handlers.Handle(serverAddr, config.ReportInterval, config.PollInterval, rep, debug)
}
