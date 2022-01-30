package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	timeout = 10
)

func HttpGet(url string, target interface{}) error {
	return HttpAuthGet(url, "", target)
}

func HttpAuthGet(url string, auth string, target interface{}) error {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if auth != "" {
		req.Header.Set("x-api-key", auth)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, filepath string) error {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	exists, err := file_exists(filepath)
	if err != nil {
		return err
	}

	var file *os.File
	if !exists {
		file, err = os.Create(filepath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
	}

	defer file.Close()
	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func file_exists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
