package main

import (
	"evo-test/internal/config"
	"evo-test/pkg/client/postgres"
	"log"
)

func main() {
	cfg := config.Get()

	log.Println("initializing postgres client...")
	_, err := postgres.NewClient(&cfg.Storage)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
		return
	}

}
