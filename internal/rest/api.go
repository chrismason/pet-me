package rest

import (
	"net/http"
	"time"

	"github.com/chrismason/pet-me/internal/config"
	"github.com/chrismason/pet-me/internal/log"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	timeout = 30
)

func timeoutMiddleware(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, time.Second*timeout, "timed out")
}

func Routes(cfg *config.ServerConfig, l *log.Logger) *mux.Router {
	mw := alice.New(l.LogMiddleware, timeoutMiddleware)

	router := mux.NewRouter()
	router.Handle("/facts/cats", mw.Then(getCatFacts(cfg, l))).Methods("GET")
	router.Handle("/pics/cats", mw.Then(getCatPic(cfg, l))).Methods("GET")
	router.Handle("/pics/dogs", mw.Then(getDogPic(cfg, l))).Methods("GET")

	return router
}
