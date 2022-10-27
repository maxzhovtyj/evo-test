package main

import (
	"evo-test/internal/config"
	"evo-test/internal/handler"
	"evo-test/internal/repository"
	"evo-test/internal/service"
	"evo-test/pkg/client/postgres"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// @title Evo Test Task
// @version 1.0
// @description API Server Evo Test Task

// @host localhost:8089
// @BasePath /
func main() {
	cfg := config.Get()

	log.Println("initializing postgres client...")
	postgresClient, err := postgres.NewClient(&cfg.Repository)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
		return
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	repoInstance := repository.New(postgresClient, psql)
	serviceInstance := service.New(repoInstance)
	handlerInstance := handler.New(serviceInstance)

	log.Println("register handlers...")
	router := handlerInstance.Register()
	log.Printf("run server on port %s...", cfg.AppPort)
	if err = runServer(router, cfg); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

func runServer(router *gin.Engine, cfg *config.Config) error {
	srv := &http.Server{
		Addr:           ":" + cfg.AppPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv.ListenAndServe()
}
