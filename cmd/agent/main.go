package main

import (
	"log"
	"metrics/internal/agent/config"
	"metrics/internal/agent/handlers"
	"metrics/internal/agent/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configs, err := config.ParseFlags()
	if err != nil {
		return err
	}

	rep := repository.NewRepository()

	return handlers.Handle(configs, rep)
}
