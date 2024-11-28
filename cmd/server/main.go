package main

import (
	"log"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	memStorage := repository.NewMemStorage()
	debug := false

	return router.Run(memStorage, debug)
}
