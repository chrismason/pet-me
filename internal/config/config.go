package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type BaseConfig struct {
	Environment string
}

type ServerConfig struct {
	BaseConfig
	ServiceName        string
	HTTPPort           int
	InstrumentationKey string
	FactsAPI           string
	CatPicsAPI         string
	CatPicsAPIKey      string
	DogPicsAPI         string
}

func Load() (*ServerConfig, error) {
	port := flag.Int("port", 8080, "port number to run the http server on")
	cfg := &ServerConfig{
		HTTPPort: *port,
	}

	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	cfg.FactsAPI = os.Getenv("CAT_FACTS_ENDPOINT")
	cfg.CatPicsAPI = os.Getenv("CAT_PICS_ENDPOINT")
	cfg.CatPicsAPIKey = os.Getenv("CAT_PICS_API_KEY")
	cfg.DogPicsAPI = os.Getenv("DOG_PICS_ENDPOINT")

	return cfg, nil
}
