package server

import (
	"net/http"

	"github.com/chrismason/pet-me/internal/config"
)

type Server struct {
	handler http.Handler
	addr    string
	cfg     *config.ServerConfig
}

func NewHTTPServer(cfg *config.ServerConfig, handler http.Handler, addr string) *Server {
	server := &Server{
		handler: handler,
		addr:    addr,
		cfg:     cfg,
	}

	return server
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.handler)
}
