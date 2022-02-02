package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chrismason/pet-me/internal/config"
	h "github.com/chrismason/pet-me/internal/http"
	"github.com/chrismason/pet-me/internal/log"
	"github.com/chrismason/pet-me/internal/models"
)

const (
	factEndpoint  = "fact"
	factsEndpoint = "facts?limit=%d"
)

func getCatFacts(cfg *config.ServerConfig, log *log.Logger) http.Handler {
	return appHandler{cfg, log, getCatFactsInner}
}

func getCatFactsInner(cfg *config.ServerConfig, log *log.Logger, w http.ResponseWriter, r *http.Request) (int, error) {
	q := r.URL.Query()
	var count int

	factCount := q.Get("factCount")

	if factCount != "" {
		var err error
		count, err = strconv.Atoi(factCount)
		if err != nil {
			e := fmt.Errorf("invalid parameter of 'factCount'")
			return http.StatusBadRequest, e
		}
		if count < 1 || count > 10 {
			return http.StatusBadRequest, fmt.Errorf("'factCount' is out of range of 1 - 10")
		}
	} else {
		count = 1
	}

	var url string
	if count > 1 {
		url = fmt.Sprintf(factsEndpoint, count)
	} else {
		url = factEndpoint
	}

	endpoint := fmt.Sprintf("%s/%s", cfg.FactsAPI, url)
	log.Info(fmt.Sprintf("Calling facts endpoint %s", endpoint))

	factsResp := models.FactResponse{}
	if count > 1 {
		facts := models.Facts{}
		err := h.HttpGet(endpoint, &facts)
		if err != nil {
			log.Dependency("cat-facts", endpoint, false)
			return http.StatusInternalServerError, fmt.Errorf("failed calling fact API")
		}
		log.Dependency("cat-facts", endpoint, true)
		for _, item := range facts.Data {
			factsResp.Facts = append(factsResp.Facts, item.Fact)
		}
	} else {
		fact := models.FactDetails{}
		err := h.HttpGet(endpoint, &fact)
		if err != nil {
			log.Dependency("cat-facts", endpoint, false)
			return http.StatusInternalServerError, fmt.Errorf("failed calling fact API")
		}
		log.Dependency("cat-facts", endpoint, true)
		factsResp.Facts = append(factsResp.Facts, fact.Fact)
	}

	err := json.NewEncoder(w).Encode(factsResp)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to encode JSON")
	}
	return http.StatusOK, nil
}
