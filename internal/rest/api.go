package rest

import (
	"net/http"
	"time"

	"github.com/chrismason/pet-me/internal/config"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	timeout = 30
)

func timeoutMiddleware(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, time.Second*timeout, "timed out")
}

func Routes(cfg *config.ServerConfig) *mux.Router {
	mw := alice.New(logMiddleware, timeoutMiddleware)

	router := mux.NewRouter()
	router.Handle("/facts/cats", mw.Then(getCatFacts(cfg))).Methods("GET")
	router.Handle("/pics/cats", mw.Then(getCatPic(cfg))).Methods("GET")
	router.Handle("/pics/dogs", mw.Then(getDogPic(cfg))).Methods("GET")

	return router
}
