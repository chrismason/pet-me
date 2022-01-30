package main

import (
	"fmt"
	"log"

	"github.com/chrismason/pet-me/internal/config"
	"github.com/chrismason/pet-me/internal/rest"
	"github.com/chrismason/pet-me/internal/server"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatalf("failed to run service: %v", err)
	}
}

func realMain() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	handler := rest.Routes(cfg)
	server := server.NewHTTPServer(cfg, handler, fmt.Sprintf(":%v", cfg.HTTPPort))

	return server.Run()
}
