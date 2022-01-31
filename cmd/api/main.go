package main

import (
	"fmt"

	"github.com/chrismason/pet-me/internal/config"
	"github.com/chrismason/pet-me/internal/log"
	"github.com/chrismason/pet-me/internal/rest"
	"github.com/chrismason/pet-me/internal/server"
)

func main() {
	if err := realMain(); err != nil {
		panic(fmt.Sprintf("failed to run service: %v", err))
	}
}

func realMain() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := log.NewLogger(log.InfoLevel, cfg)

	handler := rest.Routes(cfg, logger)
	server := server.NewHTTPServer(cfg, handler, fmt.Sprintf(":%v", cfg.HTTPPort))

	return server.Run()
}
