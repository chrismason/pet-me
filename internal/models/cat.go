package models

type CatSearchResponse struct {
	Breeds []interface{} `json:"breeds"`
	Id     string        `json:"id"`
	Url    string        `json:"url"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
}
