package main

import (
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/repository"
	"metrics/internal/server/router"
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
	memStorage := repository.NewMemStorage()

	return router.Run(configs, memStorage)
}
