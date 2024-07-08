package quoteService

import (
	"encoding/json"
	"io"
	"net/http"
)

const QuoteUrl = "https://api.chucknorris.io/jokes/random"

type QuoteHandler interface {
	GetRandom() (string, error)
}

func NewQuoteService() QuoteHandler {
	return &quoteService{url: QuoteUrl}
}

type quoteService struct {
	url string
}

func (s *quoteService) GetRandom() (string, error) {
	resp, err := http.Get(s.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	str := struct {
		Value string `json:"value"`
	}{}
	if err := json.Unmarshal(body, &str); err != nil {
		return "", err
	}
	return str.Value, nil
}
