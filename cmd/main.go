package main

import (
	"evo-test/internal/config"
	"evo-test/internal/handler"
	"evo-test/internal/repository"
	"evo-test/internal/service"
	"evo-test/pkg/client/postgres"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := config.Get()

	log.Println("initializing postgres client...")
	postgresClient, err := postgres.NewClient(&cfg.Repository)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
		return
	}

	router := httprouter.New()

	repoInstance := repository.New(postgresClient)
	serviceInstance := service.New(repoInstance)
	handlerInstance := handler.New(serviceInstance)

	log.Println("register handlers...")
	handlerInstance.Register(router)

	log.Println("run server...")
	if err = runServer(router, cfg); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

func runServer(router *httprouter.Router, cfg *config.Config) error {
	srv := &http.Server{
		Addr:           ":" + cfg.AppPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv.ListenAndServe()
}
