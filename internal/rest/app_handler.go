package rest

import (
	"net/http"

	"github.com/chrismason/pet-me/internal/config"
)

type appHandler struct {
	cfg *config.ServerConfig
	fn  func(cfg *config.ServerConfig, w http.ResponseWriter, r *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.fn(ah.cfg, w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
