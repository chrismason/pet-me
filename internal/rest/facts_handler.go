package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chrismason/pet-me/internal/config"
	h "github.com/chrismason/pet-me/internal/http"
	"github.com/chrismason/pet-me/internal/models"
)

const (
	factEndpoint  = "fact"
	factsEndpoint = "facts?limit=%d"
)

func getCatFacts(cfg *config.ServerConfig) http.Handler {
	return appHandler{cfg, getCatFactsInner}
}

func getCatFactsInner(cfg *config.ServerConfig, w http.ResponseWriter, r *http.Request) (int, error) {
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
	fmt.Printf("Calling facts endpoint %s\n", endpoint)

	factsResp := models.FactResponse{}
	if count > 1 {
		facts := models.Facts{}
		err := h.HttpGet(endpoint, &facts)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed calling fact API")
		}
		for _, item := range facts.Data {
			factsResp.Facts = append(factsResp.Facts, item.Fact)
		}
	} else {
		fact := models.FactDetails{}
		err := h.HttpGet(endpoint, &fact)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed calling fact API")
		}
		factsResp.Facts = append(factsResp.Facts, fact.Fact)
	}

	err := json.NewEncoder(w).Encode(factsResp)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to encode JSON")
	}
	return http.StatusOK, nil
}
