package adapters

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"tg_motivation_bot/internal/interfaces"
)

type TranslateAdapter struct {
	API_URL string // https://api.mymemory.translated.net
}

func NewTranslateAdapter(apiUrl string) *TranslateAdapter {
	return &TranslateAdapter{
		API_URL: apiUrl,
	}
}

// Translate(text, fromLang, toLang string) (*Quote, error)
func (t *TranslateAdapter) Translate(quote *interfaces.Quote, fromLang, toLang string) (*interfaces.Quote, error) {
	// 1. переводим цитату
	quoteText, err := t.translateText(quote.Text, fromLang, toLang)
	if err != nil {
		log.Println("Не удалось перевести цитату")
		return nil, err
	}

	// 2. переводим имя фамилию автора
	author, err := t.translateText(quote.Author, fromLang, toLang)
	if err != nil {
		log.Println("Не удалось перевести имя фамилию автора")
		return nil, err
	}

	newQuote := &interfaces.Quote{
		Text:   *quoteText,
		Author: *author,
	}

	return newQuote, nil
}

func (t *TranslateAdapter) translateText(text, fromLang, toLang string) (*string, error) {
	textQuote := strings.ReplaceAll(text, " ", "%20")
	reqUrlQuote := t.API_URL + "?q=" + textQuote + "&langpair=" + fromLang + "|" + toLang
	resp, err := http.Get(reqUrlQuote)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch translation from MyMemory")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var resultQuote MyMemoryResponse
	if err = json.Unmarshal(body, &resultQuote); err != nil {
		return nil, err
	}

	return &resultQuote.ResponseData.TranslatedText, nil
}
