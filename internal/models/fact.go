package models

type FactDetails struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type Link struct {
	Url    string `json:"url"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type Facts struct {
	CurrentPage     int           `json:"current_page"`
	Data            []FactDetails `json:"data"`
	FirstPage       string        `json:"first_page_url"`
	From            int           `json:"from"`
	LastPage        int           `json:"last_page"`
	LastPageUrl     string        `json:"last_page_url"`
	Links           []Link        `json:"links"`
	NextPageUrl     string        `json:"next_page_url"`
	Path            string        `json:"path"`
	PerPage         string        `json:"per_page"`
	PreviousPageUrl string        `json:"prev_page_url"`
	To              int           `json:"to"`
	Total           int           `json:"total"`
}

type FactResponse struct {
	Facts []string `json:"facts"`
}
