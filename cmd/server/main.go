package main

import (
	"fmt"
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
	netAddr, err := parseFlags()
	if err != nil {
		return err
	}
	memStorage := repository.NewMemStorage()
	debug := false

	serverAddr := fmt.Sprintf("%s:%d", netAddr.Host, netAddr.Port)
	return router.Run(serverAddr, memStorage, debug)
}
