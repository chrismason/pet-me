package rest

import (
	"net/http"

	"github.com/chrismason/pet-me/internal/config"
	"github.com/chrismason/pet-me/internal/log"
)

type appHandler struct {
	cfg *config.ServerConfig
	log *log.Logger
	fn  func(cfg *config.ServerConfig, log *log.Logger, w http.ResponseWriter, r *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.fn(ah.cfg, ah.log, w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
