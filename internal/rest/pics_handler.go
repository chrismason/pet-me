package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chrismason/pet-me/internal/config"
	h "github.com/chrismason/pet-me/internal/http"
	"github.com/chrismason/pet-me/internal/log"
	"github.com/chrismason/pet-me/internal/models"

	"github.com/qeesung/image2ascii/convert"
)

const (
	catSearchUrl = "images/search?mime_types=jpg,png"
	dogSearchUrl = "breeds/image/random"
	picWidth     = 100
	picHeight    = 40
)

func getCatPic(cfg *config.ServerConfig, log *log.Logger) http.Handler {
	return appHandler{cfg, log, getCatPicInner}
}

func getDogPic(cfg *config.ServerConfig, log *log.Logger) http.Handler {
	return appHandler{cfg, log, getDogPicInner}
}

func getCatPicInner(cfg *config.ServerConfig, log *log.Logger, w http.ResponseWriter, _ *http.Request) (int, error) {
	endpoint := fmt.Sprintf("%s/%s", cfg.CatPicsAPI, catSearchUrl)
	log.Info(fmt.Sprintf("Calling cat pics endpoint %s", endpoint))

	pics := []models.CatSearchResponse{}
	err := h.HttpAuthGet(endpoint, cfg.CatPicsAPIKey, &pics)

	if err != nil {
		log.Dependency("cat-pics", endpoint, false)
		return http.StatusInternalServerError, err
	}
	log.Dependency("cat-pics", endpoint, true)

	cat := pics[0]
	fileExt := filepath.Ext(cat.Url)
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("temp-*%s", fileExt))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	file.Close()
	defer os.Remove(file.Name())

	err = h.DownloadFile(cat.Url, file.Name())
	if err != nil {
		log.Dependency("cat-pic-file", cat.Url, false)
		return http.StatusInternalServerError, err
	}
	log.Dependency("cat-pic-file", cat.Url, true)

	log.Info(fmt.Sprintf("File created at %s", file.Name()))
	content := convertImageToAscii(file.Name())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	picResp := &models.Pic{
		Data:   content,
		Width:  picWidth,
		Height: picHeight,
	}

	err = json.NewEncoder(w).Encode(picResp)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func getDogPicInner(cfg *config.ServerConfig, log *log.Logger, w http.ResponseWriter, _ *http.Request) (int, error) {
	endpoint := fmt.Sprintf("%s/%s", cfg.DogPicsAPI, dogSearchUrl)
	log.Info(fmt.Sprintf("Calling dog pics endpoint %s", endpoint))

	pic := models.DogPicResponse{}
	err := h.HttpGet(endpoint, &pic)

	if err != nil {
		log.Dependency("dog-pics", endpoint, false)
		return http.StatusInternalServerError, err
	}
	log.Dependency("dog-pics", endpoint, true)

	fileExt := filepath.Ext(pic.Url)
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("temp-*%s", fileExt))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	file.Close()
	defer os.Remove(file.Name())

	err = h.DownloadFile(pic.Url, file.Name())
	if err != nil {
		log.Dependency("dog-pic-file", pic.Url, false)
		return http.StatusInternalServerError, err
	}
	log.Dependency("dog-pic-file", pic.Url, true)

	log.Info(fmt.Sprintf("File created at %s", file.Name()))
	content := convertImageToAscii(file.Name())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	fmt.Printf("%s\n", content)

	picResp := &models.Pic{
		Data:   content,
		Width:  picWidth,
		Height: picHeight,
	}

	err = json.NewEncoder(w).Encode(picResp)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func convertImageToAscii(filepath string) string {
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = picWidth
	convertOptions.FixedHeight = picHeight

	converter := convert.NewImageConverter()
	data := converter.ImageFile2ASCIIString(filepath, &convertOptions)

	return data
}
