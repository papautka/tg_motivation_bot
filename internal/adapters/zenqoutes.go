package adapters

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"tg_motivation_bot/internal/interfaces"
)

type zenQuoteResponse struct {
	Q string `json:"q"` // цитата
	A string `json:"a"` // автор
}

type ZenQuotesAdapter struct {
	API_URL string // https://zenquotes.io/api/random
}

func NewZenQuotesAdapter(apiURL string) *ZenQuotesAdapter {
	return &ZenQuotesAdapter{API_URL: apiURL}
}

func (adapter *ZenQuotesAdapter) GetQuote() (*interfaces.Quote, error) {
	// 1. делаем GET запрос
	resp, err := http.Get(adapter.API_URL)
	if err != nil {
		log.Printf("Failed to get quote from Zenqoutes: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 2. Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get quote from Zenqoutes: %s", resp.Status)
		return nil, err
	}

	// 3. Читаем тело
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Неудачная попытка чтения из цитата: %v", err)
		return nil, err
	}

	// 4. Десесериализируем данные
	var respQuote []zenQuoteResponse
	if err = json.Unmarshal(body, &respQuote); err != nil {
		log.Printf("Неудачная попытка десериализовать данные %v", err)
		return nil, err
	}

	if len(respQuote) == 0 {
		return nil, errors.New("empty quote result")
	}

	quote := &interfaces.Quote{
		Text:   respQuote[0].Q,
		Author: respQuote[0].A,
	}

	return quote, nil
}
