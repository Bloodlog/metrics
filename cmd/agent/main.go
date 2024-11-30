package main

import (
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
	rep := repository.NewRepository()
	debug := true

	return handlers.Handle(rep, debug)
}
